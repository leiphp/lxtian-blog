package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GithubClient struct {
	Config *GithubConfig
}

// NewGithubClient 创建GitHub客户端
func NewGithubClient(config *GithubConfig) *GithubClient {
	return &GithubClient{
		Config: config,
	}
}

// GetAuthURL 获取GitHub授权URL
func (c *GithubClient) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", c.Config.ClientID)
	params.Set("redirect_uri", c.Config.RedirectURL)
	params.Set("scope", strings.Join(c.Config.Scopes, " "))
	params.Set("state", state)

	return fmt.Sprintf("%s?%s", c.Config.AuthURL, params.Encode())
}

// GetAccessToken 通过code获取access_token
func (c *GithubClient) GetAccessToken(code string) (string, error) {
	params := url.Values{}
	params.Set("client_id", c.Config.ClientID)
	params.Set("client_secret", c.Config.ClientSecret)
	params.Set("code", code)
	params.Set("redirect_uri", c.Config.RedirectURL)

	req, err := http.NewRequest("POST", c.Config.TokenURL, strings.NewReader(params.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// GitHub需要Accept头指定返回JSON格式
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("获取access_token失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.Error != "" {
		return "", fmt.Errorf("获取access_token失败: %s - %s", result.Error, result.ErrorDesc)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("access_token为空: %s", string(body))
	}

	return result.AccessToken, nil
}

// GetUserInfo 获取GitHub用户信息
func (c *GithubClient) GetUserInfo(accessToken string) (*OAuthUserInfo, error) {
	req, err := http.NewRequest("GET", c.Config.UserURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// GitHub API需要Bearer token
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
		Email     string `json:"email"`
		Bio       string `json:"bio"`
		Message   string `json:"message"` // 错误信息
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	if result.Message != "" {
		return nil, fmt.Errorf("获取用户信息失败: %s", result.Message)
	}

	nickname := result.Name
	if nickname == "" {
		nickname = result.Login
	}

	// 如果主接口没有email，尝试获取email列表
	email := result.Email
	if email == "" {
		email, _ = c.getUserEmail(accessToken)
	}

	return &OAuthUserInfo{
		OpenID:      fmt.Sprintf("%d", result.ID),
		Nickname:    nickname,
		HeadImg:     result.AvatarURL,
		Email:       email,
		AccessToken: accessToken,
	}, nil
}

// getUserEmail 获取用户邮箱（GitHub API）
func (c *GithubClient) getUserEmail(accessToken string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}

	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	// 优先返回主邮箱且已验证的
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	// 其次返回已验证的
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	// 最后返回第一个
	if len(emails) > 0 {
		return emails[0].Email, nil
	}

	return "", nil
}

// RefreshToken GitHub OAuth不支持刷新token
func (c *GithubClient) RefreshToken(refreshToken string) (string, error) {
	return "", fmt.Errorf("GitHub OAuth不支持刷新token")
}

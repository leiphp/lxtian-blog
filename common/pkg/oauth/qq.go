package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type QQClient struct {
	Config *QQConfig
}

// NewQQClient 创建QQ客户端
func NewQQClient(config *QQConfig) *QQClient {
	return &QQClient{
		Config: config,
	}
}

// GetAuthURL 获取QQ授权URL
func (c *QQClient) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", c.Config.ClientID)
	params.Set("redirect_uri", c.Config.RedirectURL)
	params.Set("state", state)
	params.Set("scope", strings.Join(c.Config.Scopes, ","))

	return fmt.Sprintf("%s?%s", c.Config.AuthURL, params.Encode())
}

// GetAccessToken 通过code获取access_token
func (c *QQClient) GetAccessToken(code string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", c.Config.ClientID)
	params.Set("client_secret", c.Config.ClientSecret)
	params.Set("code", code)
	params.Set("redirect_uri", c.Config.RedirectURL)

	tokenURL := fmt.Sprintf("%s?%s", c.Config.TokenURL, params.Encode())

	resp, err := http.Get(tokenURL)
	if err != nil {
		return "", fmt.Errorf("获取access_token失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// QQ返回的是URL编码格式: access_token=xxx&expires_in=xxx
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	accessToken := values.Get("access_token")
	if accessToken == "" {
		return "", fmt.Errorf("access_token为空: %s", string(body))
	}

	return accessToken, nil
}

// GetOpenID 获取用户OpenID
func (c *QQClient) GetOpenID(accessToken string) (string, error) {
	openIDURL := fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s", accessToken)

	resp, err := http.Get(openIDURL)
	if err != nil {
		return "", fmt.Errorf("获取OpenID失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// QQ返回格式: callback( {"client_id":"YOUR_APPID","openid":"YOUR_OPENID"} );
	re := regexp.MustCompile(`"openid":"([^"]+)"`)
	matches := re.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", fmt.Errorf("无法解析OpenID: %s", string(body))
	}

	return matches[1], nil
}

// GetUserInfo 获取QQ用户信息
func (c *QQClient) GetUserInfo(accessToken string) (*OAuthUserInfo, error) {
	// 先获取OpenID
	openID, err := c.GetOpenID(accessToken)
	if err != nil {
		return nil, err
	}

	// 构建用户信息请求URL
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("oauth_consumer_key", c.Config.ClientID)
	params.Set("openid", openID)

	userInfoURL := fmt.Sprintf("%s?%s", c.Config.UserURL, params.Encode())

	resp, err := http.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Ret       int    `json:"ret"`
		Msg       string `json:"msg"`
		Nickname  string `json:"nickname"`
		Figureurl string `json:"figureurl_qq_2"` // QQ头像
		Gender    string `json:"gender"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	if result.Ret != 0 {
		return nil, fmt.Errorf("获取用户信息失败: %s", result.Msg)
	}

	// 如果没有高清头像，使用普通头像
	headImg := result.Figureurl
	if headImg == "" {
		var fallbackResult struct {
			Figureurl string `json:"figureurl_qq_1"`
		}
		json.Unmarshal(body, &fallbackResult)
		headImg = fallbackResult.Figureurl
	}

	return &OAuthUserInfo{
		OpenID:      openID,
		Nickname:    result.Nickname,
		HeadImg:     headImg,
		AccessToken: accessToken,
	}, nil
}

// RefreshToken QQ暂不支持刷新token
func (c *QQClient) RefreshToken(refreshToken string) (string, error) {
	return "", fmt.Errorf("QQ OAuth不支持刷新token")
}

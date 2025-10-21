package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type WeiboClient struct {
	Config *WeiboConfig
}

// NewWeiboClient 创建微博客户端
func NewWeiboClient(config *WeiboConfig) *WeiboClient {
	return &WeiboClient{
		Config: config,
	}
}

// GetAuthURL 获取微博授权URL
func (c *WeiboClient) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", c.Config.ClientID)
	params.Set("redirect_uri", c.Config.RedirectURL)
	params.Set("state", state)
	params.Set("scope", strings.Join(c.Config.Scopes, ","))
	params.Set("response_type", "code")

	return fmt.Sprintf("%s?%s", c.Config.AuthURL, params.Encode())
}

// GetAccessToken 通过code获取access_token
func (c *WeiboClient) GetAccessToken(code string) (string, error) {
	params := url.Values{}
	params.Set("client_id", c.Config.ClientID)
	params.Set("client_secret", c.Config.ClientSecret)
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("redirect_uri", c.Config.RedirectURL)

	resp, err := http.PostForm(c.Config.TokenURL, params)
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
		ExpiresIn   int    `json:"expires_in"`
		UID         string `json:"uid"`
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

// GetUserInfo 获取微博用户信息
func (c *WeiboClient) GetUserInfo(accessToken string) (*OAuthUserInfo, error) {
	// 先获取UID
	uidURL := fmt.Sprintf("https://api.weibo.com/2/account/get_uid.json?access_token=%s", accessToken)

	resp, err := http.Get(uidURL)
	if err != nil {
		return nil, fmt.Errorf("获取UID失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var uidResult struct {
		UID   int64  `json:"uid"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal(body, &uidResult); err != nil {
		return nil, fmt.Errorf("解析UID失败: %w", err)
	}

	if uidResult.Error != "" {
		return nil, fmt.Errorf("获取UID失败: %s", uidResult.Error)
	}

	// 获取用户信息
	userInfoURL := fmt.Sprintf("%s?access_token=%s&uid=%d", c.Config.UserURL, accessToken, uidResult.UID)

	resp, err = http.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		ID              int64  `json:"id"`
		ScreenName      string `json:"screen_name"`
		Name            string `json:"name"`
		AvatarLarge     string `json:"avatar_large"`
		AvatarHD        string `json:"avatar_hd"`
		ProfileImageURL string `json:"profile_image_url"`
		Gender          string `json:"gender"`
		Error           string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("获取用户信息失败: %s", result.Error)
	}

	// 优先使用高清头像
	headImg := result.AvatarHD
	if headImg == "" {
		headImg = result.AvatarLarge
	}
	if headImg == "" {
		headImg = result.ProfileImageURL
	}

	nickname := result.ScreenName
	if nickname == "" {
		nickname = result.Name
	}

	return &OAuthUserInfo{
		OpenID:      fmt.Sprintf("%d", result.ID),
		Nickname:    nickname,
		HeadImg:     headImg,
		AccessToken: accessToken,
	}, nil
}

// RefreshToken 刷新微博token
func (c *WeiboClient) RefreshToken(refreshToken string) (string, error) {
	return "", fmt.Errorf("微博OAuth暂不支持刷新token")
}

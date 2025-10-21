package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type WechatClient struct {
	Config *WechatConfig
}

// NewWechatClient 创建微信客户端
func NewWechatClient(config *WechatConfig) *WechatClient {
	return &WechatClient{
		Config: config,
	}
}

// GetAuthURL 获取微信授权URL（扫码登录）
func (c *WechatClient) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("appid", c.Config.ClientID)
	params.Set("redirect_uri", c.Config.RedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join(c.Config.Scopes, ","))
	params.Set("state", state)

	return fmt.Sprintf("%s?%s#wechat_redirect", c.Config.AuthURL, params.Encode())
}

// GetAccessToken 通过code获取access_token
func (c *WechatClient) GetAccessToken(code string) (string, error) {
	params := url.Values{}
	params.Set("appid", c.Config.ClientID)
	params.Set("secret", c.Config.ClientSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")

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

	var result struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
		UnionID      string `json:"unionid"`
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("获取access_token失败 [%d]: %s", result.ErrCode, result.ErrMsg)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("access_token为空: %s", string(body))
	}

	return result.AccessToken, nil
}

// GetUserInfo 获取微信用户信息
func (c *WechatClient) GetUserInfo(accessToken string) (*OAuthUserInfo, error) {
	// 先获取OpenID
	openID, unionID, err := c.getOpenID(accessToken)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("openid", openID)
	params.Set("lang", "zh_CN")

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
		OpenID     string   `json:"openid"`
		Nickname   string   `json:"nickname"`
		Sex        int      `json:"sex"`
		Province   string   `json:"province"`
		City       string   `json:"city"`
		Country    string   `json:"country"`
		HeadImgURL string   `json:"headimgurl"`
		Privilege  []string `json:"privilege"`
		UnionID    string   `json:"unionid"`
		ErrCode    int      `json:"errcode"`
		ErrMsg     string   `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %w", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("获取用户信息失败 [%d]: %s", result.ErrCode, result.ErrMsg)
	}

	// 使用返回的unionID，如果没有则使用之前获取的
	if result.UnionID != "" {
		unionID = result.UnionID
	}

	return &OAuthUserInfo{
		OpenID:      openID,
		Nickname:    result.Nickname,
		HeadImg:     result.HeadImgURL,
		UnionID:     unionID,
		AccessToken: accessToken,
	}, nil
}

// getOpenID 从access_token中获取OpenID
func (c *WechatClient) getOpenID(accessToken string) (string, string, error) {
	// 实际上在获取access_token时就已经返回了openid和unionid
	// 这里我们需要从token中提取或者重新请求
	// 简化处理：通过检查token接口获取
	checkURL := fmt.Sprintf("%s?access_token=%s&openid=", c.Config.CheckTokenURL, accessToken)

	resp, err := http.Get(checkURL)
	if err != nil {
		return "", "", fmt.Errorf("检查token失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		OpenID  string `json:"openid"`
		UnionID string `json:"unionid"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 注意：实际使用时，OpenID应该在GetAccessToken时就保存下来
	// 这里只是示例，建议在实际使用时优化这个流程
	return result.OpenID, result.UnionID, nil
}

// RefreshToken 刷新微信token
func (c *WechatClient) RefreshToken(refreshToken string) (string, error) {
	params := url.Values{}
	params.Set("appid", c.Config.ClientID)
	params.Set("grant_type", "refresh_token")
	params.Set("refresh_token", refreshToken)

	refreshURL := fmt.Sprintf("%s?%s", c.Config.RefreshURL, params.Encode())

	resp, err := http.Get(refreshURL)
	if err != nil {
		return "", fmt.Errorf("刷新token失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("刷新token失败 [%d]: %s", result.ErrCode, result.ErrMsg)
	}

	return result.AccessToken, nil
}

package oauth

// OAuthConfig OAuth配置
type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// QQConfig QQ OAuth配置
type QQConfig struct {
	OAuthConfig
	AuthURL  string
	TokenURL string
	UserURL  string
}

// WeiboConfig 微博 OAuth配置
type WeiboConfig struct {
	OAuthConfig
	AuthURL  string
	TokenURL string
	UserURL  string
}

// GithubConfig GitHub OAuth配置
type GithubConfig struct {
	OAuthConfig
	AuthURL  string
	TokenURL string
	UserURL  string
}

// WechatConfig 微信扫码 OAuth配置
type WechatConfig struct {
	OAuthConfig
	AuthURL       string
	TokenURL      string
	UserURL       string
	RefreshURL    string
	CheckTokenURL string
}

// DefaultQQConfig 默认QQ配置
func DefaultQQConfig(clientID, clientSecret, redirectURL string) *QQConfig {
	return &QQConfig{
		OAuthConfig: OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"get_user_info"},
		},
		AuthURL:  "https://graph.qq.com/oauth2.0/authorize",
		TokenURL: "https://graph.qq.com/oauth2.0/token",
		UserURL:  "https://graph.qq.com/user/get_user_info",
	}
}

// DefaultWeiboConfig 默认微博配置
func DefaultWeiboConfig(clientID, clientSecret, redirectURL string) *WeiboConfig {
	return &WeiboConfig{
		OAuthConfig: OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email"},
		},
		AuthURL:  "https://api.weibo.com/oauth2/authorize",
		TokenURL: "https://api.weibo.com/oauth2/access_token",
		UserURL:  "https://api.weibo.com/2/users/show.json",
	}
}

// DefaultGithubConfig 默认GitHub配置
func DefaultGithubConfig(clientID, clientSecret, redirectURL string) *GithubConfig {
	return &GithubConfig{
		OAuthConfig: OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user:email"},
		},
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
		UserURL:  "https://api.github.com/user",
	}
}

// DefaultWechatConfig 默认微信配置
func DefaultWechatConfig(clientID, clientSecret, redirectURL string) *WechatConfig {
	return &WechatConfig{
		OAuthConfig: OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"snsapi_login"},
		},
		AuthURL:       "https://open.weixin.qq.com/connect/qrconnect",
		TokenURL:      "https://api.weixin.qq.com/sns/oauth2/access_token",
		UserURL:       "https://api.weixin.qq.com/sns/userinfo",
		RefreshURL:    "https://api.weixin.qq.com/sns/oauth2/refresh_token",
		CheckTokenURL: "https://api.weixin.qq.com/sns/auth",
	}
}

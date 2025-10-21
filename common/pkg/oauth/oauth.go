package oauth

// OAuthUserInfo OAuth 用户信息统一接口
type OAuthUserInfo struct {
	OpenID      string `json:"openid"`       // 用户唯一标识
	Nickname    string `json:"nickname"`     // 昵称
	HeadImg     string `json:"head_img"`     // 头像
	Email       string `json:"email"`        // 邮箱（可选）
	UnionID     string `json:"unionid"`      // 联合ID（微信专用）
	AccessToken string `json:"access_token"` // 访问令牌
}

// OAuthClient OAuth客户端接口
type OAuthClient interface {
	// GetAuthURL 获取授权URL
	GetAuthURL(state string) string

	// GetAccessToken 通过code获取access_token
	GetAccessToken(code string) (string, error)

	// GetUserInfo 获取用户信息
	GetUserInfo(accessToken string) (*OAuthUserInfo, error)

	// RefreshToken 刷新token（可选）
	RefreshToken(refreshToken string) (string, error)
}

package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"lxtian-blog/common/pkg/define"
	"lxtian-blog/common/pkg/jwts"
	"lxtian-blog/common/pkg/oauth"
	"lxtian-blog/common/pkg/redis"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthCallbackLogic {
	return &AuthCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AuthCallback 通用OAuth登录-授权回调
func (l *AuthCallbackLogic) AuthCallback(w http.ResponseWriter, r *http.Request) error {
	// 从URL路径中获取OAuth类型
	oauthType := l.getOAuthTypeFromPath(r.URL.Path)
	if oauthType == "" {
		return l.redirectToFrontendWithError(w, r, "无效的OAuth类型")
	}

	// 获取授权码和state
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		return l.redirectToFrontendWithError(w, r, "授权失败：未获取到授权码")
	}

	// 验证state
	storedState, err := l.svcCtx.Rds.Get(redis.ReturnRedisKey(redis.OAuthStateString, state))
	if err != nil || storedState != state {
		logx.Errorf("state验证失败: type=%s, stored=%s, received=%s, err=%v", oauthType, storedState, state, err)
		return l.redirectToFrontendWithError(w, r, "授权失败：state验证失败")
	}

	// 删除已使用的state
	l.svcCtx.Rds.Del(redis.ReturnRedisKey(redis.OAuthStateString, state))

	// 根据类型创建对应的OAuth客户端
	oauthClient, err := l.createOAuthClient(oauthType)
	if err != nil {
		logx.Errorf("创建OAuth客户端失败: type=%s, err=%v", oauthType, err)
		return l.redirectToFrontendWithError(w, r, fmt.Sprintf("不支持的OAuth类型: %s", oauthType))
	}

	// 获取access_token
	accessToken, err := oauthClient.GetAccessToken(code)
	if err != nil {
		logx.Errorf("获取access_token失败: type=%s, err=%v", oauthType, err)
		return l.redirectToFrontendWithError(w, r, "授权失败：获取access_token失败")
	}

	// 获取用户信息
	userInfo, err := oauthClient.GetUserInfo(accessToken)
	if err != nil {
		logx.Errorf("获取用户信息失败: type=%s, err=%v", oauthType, err)
		return l.redirectToFrontendWithError(w, r, "授权失败：获取用户信息失败")
	}

	// 获取登录类型
	loginType := l.getLoginType(oauthType)

	// 调用RPC创建/更新用户（只传递用户信息，不传递code）
	// 这样RPC层不需要处理OAuth逻辑，只负责用户数据管理
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		LoginType: uint32(loginType),
		Username:  userInfo.OpenID, // OAuth平台的唯一标识
		Nickname:  userInfo.Nickname,
		HeadImg:   userInfo.HeadImg,
	})
	if err != nil {
		logx.Errorf("RPC登录失败: type=%s, err=%v", oauthType, err)
		return l.redirectToFrontendWithError(w, r, "登录失败")
	}

	var result map[string]interface{}
	if err = json.Unmarshal([]byte(res.Data), &result); err != nil {
		logx.Errorf("解析RPC响应失败: type=%s, err=%v", oauthType, err)
		return l.redirectToFrontendWithError(w, r, "登录失败")
	}

	// 生成JWT token
	auth := l.svcCtx.Config.Auth
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   uint(result["id"].(float64)),
		Username: result["username"].(string),
	}, auth.AccessSecret, auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成token失败: type=%s, err=%v", oauthType, err)
		return l.redirectToFrontendWithError(w, r, "登录失败")
	}

	// 将token存储到Redis
	err = l.svcCtx.Rds.Setex(redis.ReturnRedisKey(redis.UserTokenString, result["id"]), token, int(auth.AccessExpire)*3600)
	if err != nil {
		logx.Errorf("存储token失败: type=%s, err=%v", oauthType, err)
	}

	logx.Infof("OAuth登录成功 - 类型: %s, 用户: %v, OpenID: %s", oauthType, result["username"], userInfo.OpenID)

	// 重定向到前端，携带token
	frontendURL := l.svcCtx.Config.OAuth.FrontendURL
	redirectURL := fmt.Sprintf("%s?token=%s&expires_in=%d", frontendURL, token, auth.AccessExpire*3600)
	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}

// getOAuthTypeFromPath 从URL路径中提取OAuth类型
// 路径格式: /user/auth/:type/callback
func (l *AuthCallbackLogic) getOAuthTypeFromPath(path string) string {
	// 去掉前缀 /user/auth/ 和后缀 /callback
	path = strings.TrimPrefix(path, "/user/auth/")
	path = strings.TrimSuffix(path, "/callback")
	return strings.ToLower(path)
}

// createOAuthClient 根据类型创建对应的OAuth客户端
func (l *AuthCallbackLogic) createOAuthClient(oauthType string) (oauth.OAuthClient, error) {
	switch strings.ToLower(oauthType) {
	case "qq":
		config := oauth.DefaultQQConfig(
			l.svcCtx.Config.OAuth.QQConf.ClientID,
			l.svcCtx.Config.OAuth.QQConf.ClientSecret,
			l.svcCtx.Config.OAuth.QQConf.RedirectURL,
		)
		return oauth.NewQQClient(config), nil

	case "weibo":
		config := oauth.DefaultWeiboConfig(
			l.svcCtx.Config.OAuth.WeiboConf.AppID,
			l.svcCtx.Config.OAuth.WeiboConf.AppSecret,
			l.svcCtx.Config.OAuth.WeiboConf.RedirectURL,
		)
		return oauth.NewWeiboClient(config), nil

	case "github":
		config := oauth.DefaultGithubConfig(
			l.svcCtx.Config.OAuth.GithubConf.ClientID,
			l.svcCtx.Config.OAuth.GithubConf.ClientSecret,
			l.svcCtx.Config.OAuth.GithubConf.RedirectURL,
		)
		return oauth.NewGithubClient(config), nil

	case "wechat":
		config := oauth.DefaultWechatConfig(
			l.svcCtx.Config.OAuth.WechatConf.AppID,
			l.svcCtx.Config.OAuth.WechatConf.AppSecret,
			l.svcCtx.Config.OAuth.WechatConf.RedirectURL,
		)
		return oauth.NewWechatClient(config), nil

	default:
		return nil, fmt.Errorf("不支持的OAuth类型: %s", oauthType)
	}
}

// getLoginType 根据OAuth类型获取对应的登录类型
func (l *AuthCallbackLogic) getLoginType(oauthType string) int32 {
	switch strings.ToLower(oauthType) {
	case "qq":
		return define.QQLogin
	case "weibo":
		return define.SinaLogin
	case "github":
		return define.GithubLogin
	case "wechat":
		return define.WechatLogin
	default:
		return define.DefaultLogin
	}
}

// redirectToFrontendWithError 重定向到前端并携带错误信息
func (l *AuthCallbackLogic) redirectToFrontendWithError(w http.ResponseWriter, r *http.Request, errorMsg string) error {
	frontendURL := l.svcCtx.Config.OAuth.FrontendURL
	redirectURL := fmt.Sprintf("%s?error=%s", frontendURL, errorMsg)
	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}


package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"lxtian-blog/common/pkg/oauth"
	"lxtian-blog/common/pkg/redis"
	"lxtian-blog/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthLoginLogic {
	return &AuthLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AuthLogin 通用OAuth登录-发起授权
func (l *AuthLoginLogic) AuthLogin(w http.ResponseWriter, r *http.Request) error {
	// 从URL路径中获取OAuth类型
	oauthType := l.getOAuthTypeFromPath(r.URL.Path)
	if oauthType == "" {
		return fmt.Errorf("无效的OAuth类型")
	}

	// 生成state参数，用于防止CSRF攻击
	state := generateState()

	// 将state存储到Redis，有效期5分钟
	err := l.svcCtx.Rds.Setex(redis.ReturnRedisKey(redis.OAuthStateString, state), state, 300)
	if err != nil {
		logx.Errorf("存储state失败: %v", err)
		return fmt.Errorf("系统错误")
	}

	// 根据类型创建对应的OAuth客户端
	oauthClient, err := l.createOAuthClient(oauthType)
	if err != nil {
		logx.Errorf("创建OAuth客户端失败: type=%s, err=%v", oauthType, err)
		return fmt.Errorf("不支持的OAuth类型: %s", oauthType)
	}

	// 获取授权URL
	authURL := oauthClient.GetAuthURL(state)

	logx.Infof("OAuth登录 - 类型: %s, 重定向到: %s", oauthType, authURL)

	// 重定向到OAuth授权页面
	http.Redirect(w, r, authURL, http.StatusFound)
	return nil
}

// getOAuthTypeFromPath 从URL路径中提取OAuth类型
// 路径格式: /user/auth/:type/login
func (l *AuthLoginLogic) getOAuthTypeFromPath(path string) string {
	// 去掉前缀 /user/auth/ 和后缀 /login
	path = strings.TrimPrefix(path, "/user/auth/")
	path = strings.TrimSuffix(path, "/login")
	return strings.ToLower(path)
}

// createOAuthClient 根据类型创建对应的OAuth客户端
func (l *AuthLoginLogic) createOAuthClient(oauthType string) (oauth.OAuthClient, error) {
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

// generateState 生成随机state
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

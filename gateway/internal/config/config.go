package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct { // JWT 认证需要的密钥和过期时间配置
		AccessSecret string `json:",env=ACCESS_SECRET"`
		AccessExpire int64
	}
	RedisConfig struct {
		Host string `json:",env=REDIS_HOST"`
		Type string `json:",env=REDIS_TYPE"`
		Pass string `json:",env=REDIS_PASS"`
		Tls  bool   `json:",env=REDIS_TLS"`
	}
	WebRpc     zrpc.RpcClientConf
	UserRpc    zrpc.RpcClientConf
	PaymentRpc zrpc.RpcClientConf
	WsService  struct {
		Host string `json:",env=WS_HOST"`
		Port int
	}
	OAuth struct {
		QQConf struct {
			ClientID     string `json:",env=QQ_CLIENT_ID"`
			ClientSecret string `json:",env=QQ_CLIENT_SECRET"`
			RedirectURL  string `json:",env=QQ_REDIRECT_URL"`
		}
		FrontendURL string `json:",env=FRONTEND_URL"` // 前端地址，用于授权成功后重定向
		WeiboConf   struct {
			AppID       string `json:",env=WEIBO_APP_ID"`
			AppSecret   string `json:",env=WEIBO_APP_SECRET"`
			RedirectURL string `json:",env=WEIBO_REDIRECT_URL"`
		}

		GithubConf struct {
			ClientID     string `json:",env=GITHUB_CLIENT_ID"`
			ClientSecret string `json:",env=GITHUB_CLIENT_SECRET"`
			RedirectURL  string `json:",env=GITHUB_REDIRECT_URL"`
		}

		WechatConf struct {
			AppID       string `json:",env=WECHAT_APP_ID"`
			AppSecret   string `json:",env=WECHAT_APP_SECRET"`
			RedirectURL string `json:",env=WECHAT_REDIRECT_URL"`
		}
	}
}

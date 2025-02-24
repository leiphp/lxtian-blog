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
	WebRpc    zrpc.RpcClientConf
	UserRpc   zrpc.RpcClientConf
	WsService struct {
		Host string `json:",env=WS_HOST"`
		Port int
	}
}

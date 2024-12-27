package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		HOST     string `json:",env=DB_HOST"`
		PORT     string `json:",env=DB_PORT"`
		DATABASE string `json:",env=DB_DATABASE"`
		USERNAME string `json:",env=DB_USERNAME"`
		PASSWORD string `json:",env=DB_PASSWORD"`
	}
	RedisConfig struct {
		Host string `json:",env=REDIS_HOST"`
		Type string `json:",env=REDIS_TYPE"`
		Pass string `json:",env=REDIS_PASS"`
		Tls  bool   `json:",env=REDIS_TLS"`
	}
}

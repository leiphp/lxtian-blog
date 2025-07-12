package config

import (
	"github.com/zeromicro/go-zero/rest"
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
	Mysql struct {
		HOST     string `json:",env=DB_HOST"`
		PORT     string `json:",env=DB_PORT"`
		DATABASE string `json:",env=DB_DATABASE"`
		USERNAME string `json:",env=DB_USERNAME"`
		PASSWORD string `json:",env=DB_PASSWORD"`
	}
	QiniuOss struct {
		AccessKey string `json:",env=AccessKey"`
		SecretKey string `json:",env=SecretKey"`
		Bucket    string `json:",env=Bucket"`
		Domain    string `json:",env=Domain"`
		Region    string `json:",env=Region"`
	}
}

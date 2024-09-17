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
	MongoDB struct {
		HOST     string `json:",env=MONGODB_HOST"`
		PORT     string `json:",env=MONGODB_PORT"`
		DATABASE string `json:",env=MONGODB_DATABASE"`
		USERNAME string `json:",env=MONGODB_USERNAME"`
		PASSWORD string `json:",env=MONGODB_PASSWORD"`
	}
}

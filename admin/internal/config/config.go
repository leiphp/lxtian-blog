package config

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Etcd EtcdConf `json:",optional"` // 新增 ETCD 配置
}

type EtcdConf struct {
	discov.EtcdConf     // 内嵌官方结构体
	DialTimeout     int `json:",default=2000"` // 自定义字段（单位：毫秒）
	TTL             int `json:",default=10"`   // 自定义字段（单位：秒）
}

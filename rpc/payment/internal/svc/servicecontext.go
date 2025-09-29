package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/common/pkg/initdb"
	"lxtian-blog/rpc/payment/internal/config"
)

type ServiceContext struct {
	Config       config.Config
	DB           *gorm.DB
	Rds          *redis.Redis
	AlipayClient *alipay.AlipayClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.USERNAME,
		c.Mysql.PASSWORD,
		c.Mysql.HOST,
		c.Mysql.PORT,
		c.Mysql.DATABASE,
	)
	// 初始化GORM数据库
	mysqlDb := initdb.InitDB(dataSource)

	// 初始化Redis
	rds := initdb.InitRedis(c.RedisConfig.Host, c.RedisConfig.Type, c.RedisConfig.Pass, c.RedisConfig.Tls)

	// 初始化支付宝客户端
	alipayConfig := &alipay.AlipayConfig{
		AppId:           c.Alipay.AppId,
		AppPrivateKey:   c.Alipay.AppPrivateKey,
		AlipayPublicKey: c.Alipay.AlipayPublicKey,
		GatewayUrl:      c.Alipay.GatewayUrl,
		NotifyUrl:       c.Alipay.NotifyUrl,
		ReturnUrl:       c.Alipay.ReturnUrl,
		IsProd:          c.Alipay.IsProd,
		SignType:        c.Alipay.SignType,
		Charset:         c.Alipay.Charset,
		Format:          c.Alipay.Format,
		Version:         c.Alipay.Version,
		Timeout:         c.Alipay.Timeout,
	}

	alipayClient, err := alipay.NewAlipayClient(alipayConfig)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:       c,
		DB:           mysqlDb,
		Rds:          rds,
		AlipayClient: alipayClient,
	}
}

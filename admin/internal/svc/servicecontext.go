package svc

import (
	"fmt"
	"github.com/leiphp/gokit/pkg/sdk/qiniu"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
	"lxtian-blog/admin/internal/config"
	"lxtian-blog/admin/internal/middleware"
	"lxtian-blog/common/pkg/initdb"
)

type ServiceContext struct {
	Config        config.Config
	JwtMiddleware rest.Middleware
	Rds           *redis.Redis
	DB            *gorm.DB
	QiniuClient   *qiniu.QiniuClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := initdb.InitRedis(c.RedisConfig.Host, c.RedisConfig.Type, c.RedisConfig.Pass, c.RedisConfig.Tls)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.USERNAME,
		c.Mysql.PASSWORD,
		c.Mysql.HOST,
		c.Mysql.PORT,
		c.Mysql.DATABASE,
	)
	mysqlDb := initdb.InitDB(dataSource)
	client := qiniu.NewClient(qiniu.QiniuConfig{
		AccessKey: c.QiniuOss.AccessKey,
		SecretKey: c.QiniuOss.SecretKey,
		Bucket:    c.QiniuOss.Bucket,
		Domain:    c.QiniuOss.Domain,
		Region:    c.QiniuOss.Region,
	})
	return &ServiceContext{
		Config:        c,
		JwtMiddleware: middleware.NewJwtMiddleware(c.Auth.AccessSecret, c.Auth.AccessExpire).Handle,
		Rds:           rds,
		DB:            mysqlDb,
		QiniuClient:   client,
	}
}

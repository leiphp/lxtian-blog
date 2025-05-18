package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
	"lxtian-blog/admin/internal/config"
	"lxtian-blog/admin/internal/middleware"
	"lxtian-blog/common/pkg/initdb"
)

type ServiceContext struct {
	Config        config.Config
	JwtMiddleware rest.Middleware
	DB            *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.USERNAME,
		c.Mysql.PASSWORD,
		c.Mysql.HOST,
		c.Mysql.PORT,
		c.Mysql.DATABASE,
	)
	mysqlDb := initdb.InitDB(dataSource)
	return &ServiceContext{
		Config:        c,
		JwtMiddleware: middleware.NewJwtMiddleware(c.Auth.AccessSecret, c.Auth.AccessExpire).Handle,
		DB:            mysqlDb,
	}
}

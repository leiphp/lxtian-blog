package svc

import (
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/initdb"
	"lxtian-blog/rpc/web/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	dataSource := c.Mysql.USERNAME + ":" + c.Mysql.PASSWORD + "@tcp(" + c.Mysql.HOST + ":" + c.Mysql.PORT + ")/" + c.Mysql.DATABASE + "?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDb := initdb.InitDB(dataSource)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}

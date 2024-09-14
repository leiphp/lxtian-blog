package svc

import (
	"fmt"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/initdb"
	"lxtian-blog/rpc/web/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
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
		Config: c,
		DB:     mysqlDb,
	}
}

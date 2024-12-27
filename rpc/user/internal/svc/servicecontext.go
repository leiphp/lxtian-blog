package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/initdb"
	"lxtian-blog/rpc/user/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	//Cache  *collection.Cache
	Rds *redis.Redis
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
	//cache, err := initcache.InitCache(60*24, "user.rpc")
	//if err != nil {
	//	logc.Errorf(context.Background(), "InitCache error: %s", err)
	//}
	rds := initdb.InitRedis(c.RedisConfig.Host, c.RedisConfig.Type, c.RedisConfig.Pass, c.RedisConfig.Tls)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
		//Cache:  cache,
		Rds: rds,
	}
}

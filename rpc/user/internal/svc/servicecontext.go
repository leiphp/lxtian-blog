package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"lxtian-blog/common/pkg/initcache"
	"lxtian-blog/common/pkg/initdb"
	"lxtian-blog/rpc/user/internal/config"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Cache  *collection.Cache
	Rds    *redis.Redis
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

	// 初始化本地缓存（30 分钟过期），整个 user.rpc 进程共享
	cache, err := initcache.InitCache(30, "user.rpc")
	if err != nil {
		logx.Errorf("InitCache error: %s", err)
	}

	rds := initdb.InitRedis(c.RedisConfig.Host, c.RedisConfig.Type, c.RedisConfig.Pass, c.RedisConfig.Tls)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
		Cache:  cache,
		Rds:    rds,
	}
}

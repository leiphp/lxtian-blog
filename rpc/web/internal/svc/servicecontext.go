package svc

import (
	"fmt"
	"github.com/leiphp/gokit/pkg/sdk/qiniu"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/initdb"
	"lxtian-blog/rpc/web/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	MongoUri    string
	Rds         *redis.Redis
	QiniuClient *qiniu.QiniuClient
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
	mongoUri := initdb.InitMongoUri(c.MongoDB.USERNAME, c.MongoDB.PASSWORD, c.MongoDB.HOST, c.MongoDB.PORT)
	rds := initdb.InitRedis(c.RedisConfig.Host, c.RedisConfig.Type, c.RedisConfig.Pass, c.RedisConfig.Tls)
	client := qiniu.NewClient(qiniu.QiniuConfig{
		AccessKey: c.QiniuOss.AccessKey,
		SecretKey: c.QiniuOss.SecretKey,
		Bucket:    c.QiniuOss.Bucket,
		Domain:    c.QiniuOss.Domain,
		Region:    c.QiniuOss.Region,
	})
	return &ServiceContext{
		Config:      c,
		DB:          mysqlDb,
		MongoUri:    mongoUri,
		Rds:         rds,
		QiniuClient: client,
	}
}

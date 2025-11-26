package initdb

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/mon"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// InitDB gorm初始化
func InitDB(MysqlDataSource string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(MysqlDataSource), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用此选项后，`User` 的表将是 `user`
		},
	})
	if err != nil {
		panic("连接mysql数据库失败, error=" + err.Error())
	} else {
		fmt.Println("连接mysql数据库成功")
	}

	// 获取通用数据库对象 sql.DB ，然后配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)           // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 设置数据库最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接的最大复用时间

	return db
}

// InitMongoDB 初始化
func InitMongoDB(uri, db, collection string) *mon.Model {
	conn := mon.MustNewModel(uri, db, collection)
	return conn
}

// InitMongoUri 初始化Uri
func InitMongoUri(username, password, host, port string) string {
	// 兼容windows设置为空的写法$env:MONGODB_USERNAME=""
	if username == "${MONGODB_USERNAME}" {
		username = ""
	}
	if password == "${MONGODB_PASSWORD}" {
		password = ""
	}
	if username == "" && password == "" {
		// 没有用户名和密码的情况
		return fmt.Sprintf("mongodb://%s:%s", host, port)
	}
	// 有用户名和密码的情况
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
}

// InitRedis 初始化Redis连接
func InitRedis(host, connType, pass string, tls bool) *redis.Redis {
	// 兼容windows设置为空的写法$env:REDIS_PASS=""
	if pass == "${REDIS_PASS}" {
		pass = ""
	}
	conf := redis.RedisConf{
		Host:        host,
		Type:        connType,
		Pass:        pass,
		Tls:         tls,
		NonBlock:    false,
		PingTimeout: time.Second,
	}
	rds := redis.MustNewRedis(conf)
	return rds
}

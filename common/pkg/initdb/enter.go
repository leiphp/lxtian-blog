package initdb

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/mon"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// InitDB gorm初始化
func InitDB(MysqlDataSource string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(MysqlDataSource), &gorm.Config{})
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

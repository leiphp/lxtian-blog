package repository

import (
	"context"
	"lxtian-blog/common/pkg/initdb"
	"sync"
	"time"

	"gorm.io/gorm"
)

// RepositoryManager 仓储管理器
type RepositoryManager struct {
	db *gorm.DB

	// 事务管理器
	TransactionManager *TransactionManager

	// 单例控制
	once sync.Once
}

// NewRepositoryManager 创建仓储管理器
func NewRepositoryManager(db *gorm.DB) *RepositoryManager {
	rm := &RepositoryManager{
		db: db,
	}

	rm.initializeRepositories()
	return rm
}

// NewRepositoryManagerWithConfig 使用配置创建仓储管理器
func NewRepositoryManagerWithConfig(mysqlDataSource string) *RepositoryManager {
	db := initdb.InitDB(mysqlDataSource)
	return NewRepositoryManager(db)
}

// initializeRepositories 初始化所有仓储
func (rm *RepositoryManager) initializeRepositories() {
	rm.once.Do(func() {
		// 初始化事务管理器
		rm.TransactionManager = NewTransactionManager(rm.db)
	})
}

// GetDB 获取数据库连接
func (rm *RepositoryManager) GetDB() *gorm.DB {
	return rm.db
}

// WithTransaction 在事务中执行操作
func (rm *RepositoryManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return rm.TransactionManager.ExecuteInTransaction(ctx, fn)
}

// BatchExecute 批量执行操作（在同一个事务中）
func (rm *RepositoryManager) BatchExecute(ctx context.Context, operations []TransactionFunc) error {
	return rm.TransactionManager.BatchExecute(ctx, operations)
}

// RepositoryManagerInterface 仓储管理器接口
type RepositoryManagerInterface interface {
	// 事务管理
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	BatchExecute(ctx context.Context, operations []TransactionFunc) error

	// 数据库连接
	GetDB() *gorm.DB
}

// GlobalRepositoryManager 全局仓储管理器实例
var (
	GlobalRepositoryManager *RepositoryManager
	once                    sync.Once
)

// InitGlobalRepositoryManager 初始化全局仓储管理器
func InitGlobalRepositoryManager(db *gorm.DB) *RepositoryManager {
	once.Do(func() {
		GlobalRepositoryManager = NewRepositoryManager(db)
	})
	return GlobalRepositoryManager
}

// GetGlobalRepositoryManager 获取全局仓储管理器
func GetGlobalRepositoryManager() *RepositoryManager {
	return GlobalRepositoryManager
}

// RepositoryConfig 仓储配置
type RepositoryConfig struct {
	MysqlDataSource string // MySQL数据源
	MaxIdleConns    int    // 最大空闲连接数
	MaxOpenConns    int    // 最大打开连接数
	MaxLifetime     int    // 连接最大生存时间（秒）
}

// NewRepositoryManagerWithOptions 使用选项创建仓储管理器
func NewRepositoryManagerWithOptions(config *RepositoryConfig) *RepositoryManager {
	db := initdb.InitDB(config.MysqlDataSource)

	// 配置连接池
	if config.MaxIdleConns > 0 || config.MaxOpenConns > 0 || config.MaxLifetime > 0 {
		sqlDB, err := db.DB()
		if err == nil {
			if config.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(config.MaxIdleConns)
			}
			if config.MaxOpenConns > 0 {
				sqlDB.SetMaxOpenConns(config.MaxOpenConns)
			}
			if config.MaxLifetime > 0 {
				sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
			}
		}
	}

	return NewRepositoryManager(db)
}

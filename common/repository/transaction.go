package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// TransactionManager 事务管理器
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// ExecuteInTransaction 在事务中执行操作
func (tm *TransactionManager) ExecuteInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建包含事务的新上下文
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

// GetDB 获取数据库连接（优先使用事务）
func (tm *TransactionManager) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return tm.db.WithContext(ctx)
}

// TransactionFunc 事务执行函数类型
type TransactionFunc func(ctx context.Context) error

// TransactionOptions 事务选项
type TransactionOptions struct {
	ReadOnly bool // 是否只读事务
}

// ExecuteWithOptions 带选项的事务执行
func (tm *TransactionManager) ExecuteWithOptions(ctx context.Context, fn TransactionFunc, options *TransactionOptions) error {
	tx := tm.db.WithContext(ctx)

	return tx.Transaction(func(tx *gorm.DB) error {
		// 设置只读事务
		if options != nil && options.ReadOnly {
			tx = tx.Exec("SET TRANSACTION READ ONLY")
		}

		// 创建包含事务的新上下文
		txCtx := context.WithValue(ctx, "tx", tx)
		return fn(txCtx)
	})
}

// BatchExecute 批量执行操作（在同一个事务中）
func (tm *TransactionManager) BatchExecute(ctx context.Context, operations []TransactionFunc) error {
	return tm.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		for i, operation := range operations {
			if err := operation(txCtx); err != nil {
				return fmt.Errorf("operation %d failed: %w", i+1, err)
			}
		}
		return nil
	})
}

// TransactionalRepository 支持事务的仓储接口
type TransactionalRepository interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// TransactionalBaseRepository 支持事务的基础仓储
type TransactionalBaseRepository[T any] struct {
	*baseRepository[T]
	txManager *TransactionManager
}

// NewTransactionalBaseRepository 创建支持事务的基础仓储
func NewTransactionalBaseRepository[T any](db *gorm.DB) *TransactionalBaseRepository[T] {
	return &TransactionalBaseRepository[T]{
		baseRepository: NewBaseRepository[T](db),
		txManager:      NewTransactionManager(db),
	}
}

// WithTransaction 事务操作
func (tr *TransactionalBaseRepository[T]) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return tr.txManager.ExecuteInTransaction(ctx, fn)
}

// GetDB 获取数据库连接（优先使用事务）
func (tr *TransactionalBaseRepository[T]) GetDB(ctx context.Context) *gorm.DB {
	return tr.txManager.GetDB(ctx)
}

// TransactionalQueryBuilder 支持事务的查询构建器
type TransactionalQueryBuilder struct {
	*QueryBuilder
	txManager *TransactionManager
}

// NewTransactionalQueryBuilder 创建支持事务的查询构建器
func NewTransactionalQueryBuilder(db *gorm.DB) *TransactionalQueryBuilder {
	return &TransactionalQueryBuilder{
		QueryBuilder: NewQueryBuilder(db),
		txManager:    NewTransactionManager(db),
	}
}

// WithContext 设置上下文（支持事务）
func (tqb *TransactionalQueryBuilder) WithContext(ctx context.Context) *TransactionalQueryBuilder {
	tqb.db = tqb.txManager.GetDB(ctx)
	return tqb
}

// ExecuteInTransaction 在事务中执行查询
func (tqb *TransactionalQueryBuilder) ExecuteInTransaction(ctx context.Context, fn func(*QueryBuilder) error) error {
	return tqb.txManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		txDB := tqb.txManager.GetDB(txCtx)
		txQB := &QueryBuilder{db: txDB}
		return fn(txQB)
	})
}

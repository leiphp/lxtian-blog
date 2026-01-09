package repository

import (
	"context"
)

// BaseRepository 基础仓储接口
type BaseRepository[T any] interface {
	// 基础CRUD操作
	Create(ctx context.Context, entity *T) error
	CreateBatch(ctx context.Context, entities []*T) error
	GetByID(ctx context.Context, id uint64) (*T, error)
	GetByCondition(ctx context.Context, condition map[string]interface{}) (*T, error)
	GetList(ctx context.Context, condition map[string]interface{}, page, pageSize int, orderBy string, keywords string, fields ...string) ([]*T, int64, error)
	Update(ctx context.Context, entity *T) error
	UpdateByCondition(ctx context.Context, condition map[string]interface{}, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint64) error
	DeleteByCondition(ctx context.Context, condition map[string]interface{}) error
	ForceDelete(ctx context.Context, id uint64) error
	ForceDeleteByCondition(ctx context.Context, condition map[string]interface{}) error
	Count(ctx context.Context, condition map[string]interface{}) (int64, error)
	Exists(ctx context.Context, condition map[string]interface{}) (bool, error)

	// 事务操作
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// 注意：具体的仓储接口已移动到各自的模块文件中
// 这里只保留基础接口定义

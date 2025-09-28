package web

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository"

	"gorm.io/gorm"
)

// TxyCategoryRepository TxyCategory表仓储接口
type TxyCategoryRepository interface {
	repository.BaseRepository[mysql.TxyCategory]

	// 分类特有方法
	GetBySlug(ctx context.Context, slug string) (*mysql.TxyCategory, error)
	GetByStatus(ctx context.Context, status int64) ([]*mysql.TxyCategory, error)
	GetByParentId(ctx context.Context, parentId uint64) ([]*mysql.TxyCategory, error)
	GetTree(ctx context.Context) ([]*mysql.TxyCategory, error)
	GetActiveCategories(ctx context.Context) ([]*mysql.TxyCategory, error)

	// 更新方法
	UpdateStatus(ctx context.Context, categoryId uint64, status int64) error
	UpdateSortOrder(ctx context.Context, categoryId uint64, sortOrder int64) error
	UpdateParentId(ctx context.Context, categoryId uint64, parentId uint64) error

	// 统计方法
	GetCountByStatus(ctx context.Context, status int64) (int64, error)
	GetCountByParentId(ctx context.Context, parentId uint64) (int64, error)
	GetArticleCountByCategory(ctx context.Context, categoryId uint64) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, categoryIds []uint64, status int64) error
	BatchUpdateSortOrder(ctx context.Context, categoryIds []uint64, sortOrder int64) error
}

// txyCategoryRepository TxyCategory表仓储实现
type txyCategoryRepository struct {
	*repository.TransactionalBaseRepository[mysql.TxyCategory]
}

// NewTxyCategoryRepository 创建TxyCategory仓储
func NewTxyCategoryRepository(db *gorm.DB) TxyCategoryRepository {
	return &txyCategoryRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[mysql.TxyCategory](db),
	}
}

// GetBySlug 根据slug获取分类
func (r *txyCategoryRepository) GetBySlug(ctx context.Context, slug string) (*mysql.TxyCategory, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"slug": slug,
	})
}

// GetByStatus 根据状态获取分类列表
func (r *txyCategoryRepository) GetByStatus(ctx context.Context, status int64) ([]*mysql.TxyCategory, error) {
	categories, _, err := r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, 0, 0) // 不分页
	return categories, err
}

// GetByParentId 根据父级ID获取分类列表
func (r *txyCategoryRepository) GetByParentId(ctx context.Context, parentId uint64) ([]*mysql.TxyCategory, error) {
	categories, _, err := r.GetList(ctx, map[string]interface{}{
		"parent_id": parentId,
	}, 0, 0) // 不分页
	return categories, err
}

// GetTree 获取分类树
func (r *txyCategoryRepository) GetTree(ctx context.Context) ([]*mysql.TxyCategory, error) {
	db := r.GetDB(ctx)
	var categories []*mysql.TxyCategory

	err := db.Where("status = ?", 1). // 启用状态
						Order("sort_order ASC, id ASC").
						Find(&categories).Error

	return categories, err
}

// GetActiveCategories 获取启用的分类列表
func (r *txyCategoryRepository) GetActiveCategories(ctx context.Context) ([]*mysql.TxyCategory, error) {
	return r.GetByStatus(ctx, 1) // 1表示启用状态
}

// UpdateStatus 更新分类状态
func (r *txyCategoryRepository) UpdateStatus(ctx context.Context, categoryId uint64, status int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": categoryId},
		map[string]interface{}{"status": status},
	)
}

// UpdateSortOrder 更新排序
func (r *txyCategoryRepository) UpdateSortOrder(ctx context.Context, categoryId uint64, sortOrder int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": categoryId},
		map[string]interface{}{"sort_order": sortOrder},
	)
}

// UpdateParentId 更新父级ID
func (r *txyCategoryRepository) UpdateParentId(ctx context.Context, categoryId uint64, parentId uint64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": categoryId},
		map[string]interface{}{"parent_id": parentId},
	)
}

// GetCountByStatus 根据状态统计分类数量
func (r *txyCategoryRepository) GetCountByStatus(ctx context.Context, status int64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"status": status,
	})
}

// GetCountByParentId 根据父级ID统计分类数量
func (r *txyCategoryRepository) GetCountByParentId(ctx context.Context, parentId uint64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"parent_id": parentId,
	})
}

// GetArticleCountByCategory 根据分类统计文章数量
func (r *txyCategoryRepository) GetArticleCountByCategory(ctx context.Context, categoryId uint64) (int64, error) {
	db := r.GetDB(ctx)
	var count int64

	err := db.Model(&mysql.TxyArticle{}).
		Where("category_id = ? AND status = ?", categoryId, 1).
		Count(&count).Error

	return count, err
}

// BatchUpdateStatus 批量更新状态
func (r *txyCategoryRepository) BatchUpdateStatus(ctx context.Context, categoryIds []uint64, status int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyCategory{}).
		Where("id IN ?", categoryIds).
		Update("status", status).Error
}

// BatchUpdateSortOrder 批量更新排序
func (r *txyCategoryRepository) BatchUpdateSortOrder(ctx context.Context, categoryIds []uint64, sortOrder int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyCategory{}).
		Where("id IN ?", categoryIds).
		Update("sort_order", sortOrder).Error
}

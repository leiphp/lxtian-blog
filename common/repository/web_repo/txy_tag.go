package web_repo

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository"

	"gorm.io/gorm"
)

// TxyTagRepository TxyTag表仓储接口
type TxyTagRepository interface {
	repository.BaseRepository[mysql.TxyTag]

	// 标签特有方法
	GetBySlug(ctx context.Context, slug string) (*mysql.TxyTag, error)
	GetByStatus(ctx context.Context, status int64) ([]*mysql.TxyTag, error)
	GetPopularTags(ctx context.Context, limit int) ([]*mysql.TxyTag, error)
	GetTagsByArticleId(ctx context.Context, articleId uint64) ([]*mysql.TxyTag, error)
	GetActiveTags(ctx context.Context) ([]*mysql.TxyTag, error)

	// 更新方法
	UpdateStatus(ctx context.Context, tagId uint64, status int64) error
	UpdateArticleCount(ctx context.Context, tagId uint64, count int64) error
	IncrementArticleCount(ctx context.Context, tagId uint64) error
	DecrementArticleCount(ctx context.Context, tagId uint64) error

	// 统计方法
	GetCountByStatus(ctx context.Context, status int64) (int64, error)
	GetArticleCountByTag(ctx context.Context, tagId uint64) (int64, error)
	GetTotalArticleCount(ctx context.Context) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, tagIds []uint64, status int64) error
	BatchUpdateArticleCount(ctx context.Context, tagIds []uint64, count int64) error
}

// txyTagRepository TxyTag表仓储实现
type txyTagRepository struct {
	*repository.TransactionalBaseRepository[mysql.TxyTag]
}

// NewTxyTagRepository 创建TxyTag仓储
func NewTxyTagRepository(db *gorm.DB) TxyTagRepository {
	return &txyTagRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[mysql.TxyTag](db),
	}
}

// GetBySlug 根据slug获取标签
func (r *txyTagRepository) GetBySlug(ctx context.Context, slug string) (*mysql.TxyTag, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"slug": slug,
	})
}

// GetByStatus 根据状态获取标签列表
func (r *txyTagRepository) GetByStatus(ctx context.Context, status int64) ([]*mysql.TxyTag, error) {
	tags, _, err := r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, 0, 0, "", "") // 不分页
	return tags, err
}

// GetPopularTags 获取热门标签
func (r *txyTagRepository) GetPopularTags(ctx context.Context, limit int) ([]*mysql.TxyTag, error) {
	db := r.GetDB(ctx)
	var tags []*mysql.TxyTag

	err := db.Where("status = ?", 1). // 启用状态
						Order("article_count DESC").
						Limit(limit).
						Find(&tags).Error

	return tags, err
}

// GetTagsByArticleId 根据文章ID获取标签列表
func (r *txyTagRepository) GetTagsByArticleId(ctx context.Context, articleId uint64) ([]*mysql.TxyTag, error) {
	db := r.GetDB(ctx)
	var tags []*mysql.TxyTag

	err := db.Table("txy_tag t").
		Select("t.*").
		Joins("LEFT JOIN txy_article_tag at ON t.id = at.tag_id").
		Where("at.article_id = ?", articleId).
		Find(&tags).Error

	return tags, err
}

// GetActiveTags 获取启用的标签列表
func (r *txyTagRepository) GetActiveTags(ctx context.Context) ([]*mysql.TxyTag, error) {
	return r.GetByStatus(ctx, 1) // 1表示启用状态
}

// UpdateStatus 更新标签状态
func (r *txyTagRepository) UpdateStatus(ctx context.Context, tagId uint64, status int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": tagId},
		map[string]interface{}{"status": status},
	)
}

// UpdateArticleCount 更新文章数量
func (r *txyTagRepository) UpdateArticleCount(ctx context.Context, tagId uint64, count int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": tagId},
		map[string]interface{}{"article_count": count},
	)
}

// IncrementArticleCount 增加文章数量
func (r *txyTagRepository) IncrementArticleCount(ctx context.Context, tagId uint64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyTag{}).
		Where("id = ?", tagId).
		Update("article_count", gorm.Expr("article_count + ?", 1)).Error
}

// DecrementArticleCount 减少文章数量
func (r *txyTagRepository) DecrementArticleCount(ctx context.Context, tagId uint64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyTag{}).
		Where("id = ?", tagId).
		Update("article_count", gorm.Expr("article_count - ?", 1)).Error
}

// GetCountByStatus 根据状态统计标签数量
func (r *txyTagRepository) GetCountByStatus(ctx context.Context, status int64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"status": status,
	})
}

// GetArticleCountByTag 根据标签统计文章数量
func (r *txyTagRepository) GetArticleCountByTag(ctx context.Context, tagId uint64) (int64, error) {
	db := r.GetDB(ctx)
	var count int64

	err := db.Table("txy_article a").
		Joins("LEFT JOIN txy_article_tag at ON a.id = at.article_id").
		Where("at.tag_id = ? AND a.status = ?", tagId, 1).
		Count(&count).Error

	return count, err
}

// GetTotalArticleCount 获取总文章数量
func (r *txyTagRepository) GetTotalArticleCount(ctx context.Context) (int64, error) {
	db := r.GetDB(ctx)
	var total int64

	err := db.Model(&mysql.TxyTag{}).
		Select("COALESCE(SUM(article_count), 0)").
		Scan(&total).Error

	return total, err
}

// BatchUpdateStatus 批量更新状态
func (r *txyTagRepository) BatchUpdateStatus(ctx context.Context, tagIds []uint64, status int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyTag{}).
		Where("id IN ?", tagIds).
		Update("status", status).Error
}

// BatchUpdateArticleCount 批量更新文章数量
func (r *txyTagRepository) BatchUpdateArticleCount(ctx context.Context, tagIds []uint64, count int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyTag{}).
		Where("id IN ?", tagIds).
		Update("article_count", count).Error
}

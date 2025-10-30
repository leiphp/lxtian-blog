package web_repo

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository"
	"time"

	"gorm.io/gorm"
)

// TxyArticleRepository TxyArticle表仓储接口
type TxyArticleRepository interface {
	repository.BaseRepository[mysql.TxyArticle]

	// 文章特有方法
	GetByTitle(ctx context.Context, title string) (*mysql.TxyArticle, error)
	GetByAuthorId(ctx context.Context, authorId uint64, page, pageSize int) ([]*mysql.TxyArticle, int64, error)
	GetByCategoryId(ctx context.Context, categoryId uint64, page, pageSize int) ([]*mysql.TxyArticle, int64, error)
	GetByStatus(ctx context.Context, status int64, page, pageSize int) ([]*mysql.TxyArticle, int64, error)
	GetByTagId(ctx context.Context, tagId uint64, page, pageSize int) ([]*mysql.TxyArticle, int64, error)
	GetPublishedArticles(ctx context.Context, page, pageSize int) ([]*mysql.TxyArticle, int64, error)
	GetPopularArticles(ctx context.Context, limit int) ([]*mysql.TxyArticle, error)
	GetLatestArticles(ctx context.Context, limit int) ([]*mysql.TxyArticle, error)
	SearchArticles(ctx context.Context, keyword string, page, pageSize int) ([]*mysql.TxyArticle, int64, error)

	// 更新方法
	UpdateStatus(ctx context.Context, articleId uint64, status int64) error
	UpdateViewCount(ctx context.Context, articleId uint64) error
	UpdateLikeCount(ctx context.Context, articleId uint64, increment int) error
	UpdateCommentCount(ctx context.Context, articleId uint64, increment int) error

	// 统计方法
	GetCountByAuthorId(ctx context.Context, authorId uint64) (int64, error)
	GetCountByCategoryId(ctx context.Context, categoryId uint64) (int64, error)
	GetCountByStatus(ctx context.Context, status int64) (int64, error)
	GetTotalViewCount(ctx context.Context) (int64, error)
	GetTotalLikeCount(ctx context.Context) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, articleIds []uint64, status int64) error
	BatchDelete(ctx context.Context, articleIds []uint64) error
	GetExpiredArticles(ctx context.Context, days int) ([]*mysql.TxyArticle, error)
}

// txyArticleRepository TxyArticle表仓储实现
type txyArticleRepository struct {
	*repository.TransactionalBaseRepository[mysql.TxyArticle]
}

// NewTxyArticleRepository 创建TxyArticle仓储
func NewTxyArticleRepository(db *gorm.DB) TxyArticleRepository {
	return &txyArticleRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[mysql.TxyArticle](db),
	}
}

// GetByTitle 根据标题获取文章
func (r *txyArticleRepository) GetByTitle(ctx context.Context, title string) (*mysql.TxyArticle, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"title": title,
	})
}

// GetByAuthorId 根据作者ID获取文章列表
func (r *txyArticleRepository) GetByAuthorId(ctx context.Context, authorId uint64, page, pageSize int) ([]*mysql.TxyArticle, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"author_id": authorId,
	}, page, pageSize, "", "")
}

// GetByCategoryId 根据分类ID获取文章列表
func (r *txyArticleRepository) GetByCategoryId(ctx context.Context, categoryId uint64, page, pageSize int) ([]*mysql.TxyArticle, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"category_id": categoryId,
	}, page, pageSize, "", "")
}

// GetByStatus 根据状态获取文章列表
func (r *txyArticleRepository) GetByStatus(ctx context.Context, status int64, page, pageSize int) ([]*mysql.TxyArticle, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, page, pageSize, "", "")
}

// GetByTagId 根据标签ID获取文章列表
func (r *txyArticleRepository) GetByTagId(ctx context.Context, tagId uint64, page, pageSize int) ([]*mysql.TxyArticle, int64, error) {
	db := r.GetDB(ctx)
	var articles []*mysql.TxyArticle
	var total int64

	// 通过文章标签关联表查询
	query := db.Table("txy_article a").
		Select("a.*").
		Joins("LEFT JOIN txy_article_tag at ON a.id = at.article_id").
		Where("at.tag_id = ?", tagId)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetPublishedArticles 获取已发布文章
func (r *txyArticleRepository) GetPublishedArticles(ctx context.Context, page, pageSize int) ([]*mysql.TxyArticle, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"status": 1, // 假设1表示已发布
	}, page, pageSize, "", "")
}

// GetPopularArticles 获取热门文章
func (r *txyArticleRepository) GetPopularArticles(ctx context.Context, limit int) ([]*mysql.TxyArticle, error) {
	db := r.GetDB(ctx)
	var articles []*mysql.TxyArticle

	err := db.Where("status = ?", 1). // 已发布
						Order("view_count DESC").
						Limit(limit).
						Find(&articles).Error

	return articles, err
}

// GetLatestArticles 获取最新文章
func (r *txyArticleRepository) GetLatestArticles(ctx context.Context, limit int) ([]*mysql.TxyArticle, error) {
	db := r.GetDB(ctx)
	var articles []*mysql.TxyArticle

	err := db.Where("status = ?", 1). // 已发布
						Order("created_at DESC").
						Limit(limit).
						Find(&articles).Error

	return articles, err
}

// SearchArticles 搜索文章
func (r *txyArticleRepository) SearchArticles(ctx context.Context, keyword string, page, pageSize int) ([]*mysql.TxyArticle, int64, error) {
	db := r.GetDB(ctx)
	var articles []*mysql.TxyArticle
	var total int64

	// 构建搜索条件
	searchCondition := "%" + keyword + "%"
	query := db.Where("status = ? AND (title LIKE ? OR content LIKE ? OR summary LIKE ?)",
		1, searchCondition, searchCondition, searchCondition)

	// 获取总数
	if err := query.Model(&mysql.TxyArticle{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// UpdateStatus 更新文章状态
func (r *txyArticleRepository) UpdateStatus(ctx context.Context, articleId uint64, status int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": articleId},
		map[string]interface{}{"status": status},
	)
}

// UpdateViewCount 更新浏览量
func (r *txyArticleRepository) UpdateViewCount(ctx context.Context, articleId uint64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyArticle{}).
		Where("id = ?", articleId).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// UpdateLikeCount 更新点赞数
func (r *txyArticleRepository) UpdateLikeCount(ctx context.Context, articleId uint64, increment int) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyArticle{}).
		Where("id = ?", articleId).
		Update("like_count", gorm.Expr("like_count + ?", increment)).Error
}

// UpdateCommentCount 更新评论数
func (r *txyArticleRepository) UpdateCommentCount(ctx context.Context, articleId uint64, increment int) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyArticle{}).
		Where("id = ?", articleId).
		Update("comment_count", gorm.Expr("comment_count + ?", increment)).Error
}

// GetCountByAuthorId 根据作者ID统计文章数量
func (r *txyArticleRepository) GetCountByAuthorId(ctx context.Context, authorId uint64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"author_id": authorId,
	})
}

// GetCountByCategoryId 根据分类ID统计文章数量
func (r *txyArticleRepository) GetCountByCategoryId(ctx context.Context, categoryId uint64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"category_id": categoryId,
	})
}

// GetCountByStatus 根据状态统计文章数量
func (r *txyArticleRepository) GetCountByStatus(ctx context.Context, status int64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"status": status,
	})
}

// GetTotalViewCount 获取总浏览量
func (r *txyArticleRepository) GetTotalViewCount(ctx context.Context) (int64, error) {
	db := r.GetDB(ctx)
	var total int64

	err := db.Model(&mysql.TxyArticle{}).
		Select("COALESCE(SUM(view_count), 0)").
		Scan(&total).Error

	return total, err
}

// GetTotalLikeCount 获取总点赞数
func (r *txyArticleRepository) GetTotalLikeCount(ctx context.Context) (int64, error) {
	db := r.GetDB(ctx)
	var total int64

	err := db.Model(&mysql.TxyArticle{}).
		Select("COALESCE(SUM(like_count), 0)").
		Scan(&total).Error

	return total, err
}

// BatchUpdateStatus 批量更新状态
func (r *txyArticleRepository) BatchUpdateStatus(ctx context.Context, articleIds []uint64, status int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyArticle{}).
		Where("id IN ?", articleIds).
		Update("status", status).Error
}

// BatchDelete 批量删除
func (r *txyArticleRepository) BatchDelete(ctx context.Context, articleIds []uint64) error {
	db := r.GetDB(ctx)
	return db.Where("id IN ?", articleIds).Delete(&mysql.TxyArticle{}).Error
}

// GetExpiredArticles 获取过期文章
func (r *txyArticleRepository) GetExpiredArticles(ctx context.Context, days int) ([]*mysql.TxyArticle, error) {
	db := r.GetDB(ctx)
	var articles []*mysql.TxyArticle

	// 计算过期时间戳
	cutoffTime := time.Now().Unix() - int64(days*24*3600)

	err := db.Where("created_at < ? AND status = ?", cutoffTime, 0). // 0表示草稿状态
										Order("created_at ASC").
										Find(&articles).Error

	return articles, err
}

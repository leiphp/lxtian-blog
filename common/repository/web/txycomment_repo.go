package web

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository"
	"time"

	"gorm.io/gorm"
)

// TxyCommentRepository TxyComment表仓储接口
type TxyCommentRepository interface {
	repository.BaseRepository[mysql.TxyComment]

	// 评论特有方法
	GetByArticleId(ctx context.Context, articleId uint64, page, pageSize int) ([]*mysql.TxyComment, int64, error)
	GetByUserId(ctx context.Context, userId uint64, page, pageSize int) ([]*mysql.TxyComment, int64, error)
	GetByParentId(ctx context.Context, parentId uint64) ([]*mysql.TxyComment, error)
	GetByStatus(ctx context.Context, status int64, page, pageSize int) ([]*mysql.TxyComment, int64, error)
	GetRecentComments(ctx context.Context, limit int) ([]*mysql.TxyComment, error)
	GetCommentsByTimeRange(ctx context.Context, startTime, endTime int64, page, pageSize int) ([]*mysql.TxyComment, int64, error)

	// 更新方法
	UpdateStatus(ctx context.Context, commentId uint64, status int64) error
	UpdateLikeCount(ctx context.Context, commentId uint64, increment int) error
	UpdateReplyCount(ctx context.Context, commentId uint64, increment int) error

	// 统计方法
	GetCountByArticleId(ctx context.Context, articleId uint64) (int64, error)
	GetCountByUserId(ctx context.Context, userId uint64) (int64, error)
	GetCountByStatus(ctx context.Context, status int64) (int64, error)
	GetTotalLikeCount(ctx context.Context) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, commentIds []uint64, status int64) error
	BatchDelete(ctx context.Context, commentIds []uint64) error
	GetExpiredComments(ctx context.Context, days int) ([]*mysql.TxyComment, error)
}

// txyCommentRepository TxyComment表仓储实现
type txyCommentRepository struct {
	*repository.TransactionalBaseRepository[mysql.TxyComment]
}

// NewTxyCommentRepository 创建TxyComment仓储
func NewTxyCommentRepository(db *gorm.DB) TxyCommentRepository {
	return &txyCommentRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[mysql.TxyComment](db),
	}
}

// GetByArticleId 根据文章ID获取评论列表
func (r *txyCommentRepository) GetByArticleId(ctx context.Context, articleId uint64, page, pageSize int) ([]*mysql.TxyComment, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"article_id": articleId,
	}, page, pageSize)
}

// GetByUserId 根据用户ID获取评论列表
func (r *txyCommentRepository) GetByUserId(ctx context.Context, userId uint64, page, pageSize int) ([]*mysql.TxyComment, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"user_id": userId,
	}, page, pageSize)
}

// GetByParentId 根据父评论ID获取子评论列表
func (r *txyCommentRepository) GetByParentId(ctx context.Context, parentId uint64) ([]*mysql.TxyComment, error) {
	comments, _, err := r.GetList(ctx, map[string]interface{}{
		"parent_id": parentId,
	}, 0, 0) // 不分页
	return comments, err
}

// GetByStatus 根据状态获取评论列表
func (r *txyCommentRepository) GetByStatus(ctx context.Context, status int64, page, pageSize int) ([]*mysql.TxyComment, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, page, pageSize)
}

// GetRecentComments 获取最近评论
func (r *txyCommentRepository) GetRecentComments(ctx context.Context, limit int) ([]*mysql.TxyComment, error) {
	db := r.GetDB(ctx)
	var comments []*mysql.TxyComment

	err := db.Where("status = ?", 1). // 已审核
						Order("created_at DESC").
						Limit(limit).
						Find(&comments).Error

	return comments, err
}

// GetCommentsByTimeRange 根据时间范围获取评论
func (r *txyCommentRepository) GetCommentsByTimeRange(ctx context.Context, startTime, endTime int64, page, pageSize int) ([]*mysql.TxyComment, int64, error) {
	db := r.GetDB(ctx)
	var comments []*mysql.TxyComment
	var total int64

	query := db.Where("created_at BETWEEN ? AND ?", startTime, endTime)

	// 获取总数
	if err := query.Model(&mysql.TxyComment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// UpdateStatus 更新评论状态
func (r *txyCommentRepository) UpdateStatus(ctx context.Context, commentId uint64, status int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": commentId},
		map[string]interface{}{"status": status},
	)
}

// UpdateLikeCount 更新点赞数
func (r *txyCommentRepository) UpdateLikeCount(ctx context.Context, commentId uint64, increment int) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyComment{}).
		Where("id = ?", commentId).
		Update("like_count", gorm.Expr("like_count + ?", increment)).Error
}

// UpdateReplyCount 更新回复数
func (r *txyCommentRepository) UpdateReplyCount(ctx context.Context, commentId uint64, increment int) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyComment{}).
		Where("id = ?", commentId).
		Update("reply_count", gorm.Expr("reply_count + ?", increment)).Error
}

// GetCountByArticleId 根据文章ID统计评论数量
func (r *txyCommentRepository) GetCountByArticleId(ctx context.Context, articleId uint64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"article_id": articleId,
	})
}

// GetCountByUserId 根据用户ID统计评论数量
func (r *txyCommentRepository) GetCountByUserId(ctx context.Context, userId uint64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"user_id": userId,
	})
}

// GetCountByStatus 根据状态统计评论数量
func (r *txyCommentRepository) GetCountByStatus(ctx context.Context, status int64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"status": status,
	})
}

// GetTotalLikeCount 获取总点赞数
func (r *txyCommentRepository) GetTotalLikeCount(ctx context.Context) (int64, error) {
	db := r.GetDB(ctx)
	var total int64

	err := db.Model(&mysql.TxyComment{}).
		Select("COALESCE(SUM(like_count), 0)").
		Scan(&total).Error

	return total, err
}

// BatchUpdateStatus 批量更新状态
func (r *txyCommentRepository) BatchUpdateStatus(ctx context.Context, commentIds []uint64, status int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyComment{}).
		Where("id IN ?", commentIds).
		Update("status", status).Error
}

// BatchDelete 批量删除
func (r *txyCommentRepository) BatchDelete(ctx context.Context, commentIds []uint64) error {
	db := r.GetDB(ctx)
	return db.Where("id IN ?", commentIds).Delete(&mysql.TxyComment{}).Error
}

// GetExpiredComments 获取过期评论
func (r *txyCommentRepository) GetExpiredComments(ctx context.Context, days int) ([]*mysql.TxyComment, error) {
	db := r.GetDB(ctx)
	var comments []*mysql.TxyComment

	// 计算过期时间戳
	cutoffTime := time.Now().Unix() - int64(days*24*3600)

	err := db.Where("created_at < ? AND status = ?", cutoffTime, 0). // 0表示待审核状态
										Order("created_at ASC").
										Find(&comments).Error

	return comments, err
}

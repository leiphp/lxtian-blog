package web_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"lxtian-blog/common/model"
	redisutil "lxtian-blog/common/pkg/redis"
	"lxtian-blog/common/repository"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type TxyDocsRepository interface {
	repository.BaseRepository[model.TxyDoc]

	// IncrementDocView 更新文档浏览量（view字段+1），一个IP每天只能算一次
	IncrementDocView(ctx context.Context, docID int32, clientIP string, rds *redis.Redis) error
	// GetDocDetail 获取文档详情，先从Redis缓存获取，如果没有则从数据库查询并加入缓存
	GetDocDetail(ctx context.Context, docID int32, rds *redis.Redis) (*model.TxyDoc, error)
}

type txyDocsRepository struct {
	*repository.TransactionalBaseRepository[model.TxyDoc]
}

func NewTxyDocsRepository(db *gorm.DB) TxyDocsRepository {
	return &txyDocsRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.TxyDoc](db),
	}
}

// IncrementDocView 更新文档浏览量（view字段+1），一个IP每天只能算一次
func (r *txyDocsRepository) IncrementDocView(ctx context.Context, docID int32, clientIP string, rds *redis.Redis) error {
	// 使用统一的Redis Key管理
	redisKey := redisutil.GetDocViewKeyToday(docID, clientIP)

	// 检查今天是否已经访问过
	exists, err := rds.ExistsCtx(ctx, redisKey)
	if err != nil {
		logc.Errorf(ctx, "检查Redis key失败: %s", err)
		return err
	}

	// 如果今天已经访问过，直接返回
	if exists {
		logc.Infof(ctx, "IP %s 今天已经访问过文档 %d", clientIP, docID)
		return nil
	}

	// 设置Redis key，过期时间为明天凌晨
	expireTime := getTomorrowMidnight()
	err = rds.SetexCtx(ctx, redisKey, "1", int(expireTime.Seconds()))
	if err != nil {
		logc.Errorf(ctx, "设置Redis key失败: %s", err)
		return err
	}

	// 更新数据库中的浏览次数
	err = r.updateDocViewCount(ctx, docID)
	if err != nil {
		logc.Errorf(ctx, "更新文档浏览次数失败: %s", err)
		return err
	}

	logc.Infof(ctx, "成功记录IP %s 访问文档 %d", clientIP, docID)
	return nil
}

// updateDocViewCount 更新数据库中的文档浏览次数
func (r *txyDocsRepository) updateDocViewCount(ctx context.Context, docID int32) error {
	// 使用原子操作更新浏览次数
	db := r.GetDB(ctx)
	err := db.WithContext(ctx).
		Table("txy_docs").
		Where("id = ?", docID).
		Update("view", gorm.Expr("view + 1")).Error

	if err != nil {
		return fmt.Errorf("更新文档浏览次数失败: %w", err)
	}

	return nil
}

// getTomorrowMidnight 获取明天凌晨的时间
func getTomorrowMidnight() time.Duration {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	tomorrowMidnight := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	return tomorrowMidnight.Sub(now)
}

// GetDocDetail 获取文档详情，先从Redis缓存获取，如果没有则从数据库查询并加入缓存
func (r *txyDocsRepository) GetDocDetail(ctx context.Context, docID int32, rds *redis.Redis) (*model.TxyDoc, error) {
	// 1. 尝试从缓存获取
	cacheKey := redisutil.ReturnRedisKey(redisutil.ApiWebStringDocDetail, docID)
	cachedData, err := rds.GetCtx(ctx, cacheKey)
	if err == nil && cachedData != "" {
		// 缓存命中，解析JSON并返回
		var doc model.TxyDoc
		if err := json.Unmarshal([]byte(cachedData), &doc); err == nil {
			logc.Infof(ctx, "从缓存获取文档详情: %d", docID)
			return &doc, nil
		}
		// JSON解析失败，继续从数据库查询
		logc.Errorf(ctx, "缓存数据解析失败，从数据库查询: %d, error: %s", docID, err)
	}

	// 2. 从数据库查询
	db := r.GetDB(ctx)
	var doc model.TxyDoc
	err = db.WithContext(ctx).Where("id = ?", docID).First(&doc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文档不存在: id=%d", docID)
		}
		return nil, fmt.Errorf("查询文档失败: %w", err)
	}

	// 3. 将查询结果加入缓存
	jsonData, err := json.Marshal(doc)
	if err == nil {
		// 设置1小时过期
		err = rds.SetexCtx(ctx, cacheKey, string(jsonData), 3600)
		if err != nil {
			logc.Errorf(ctx, "设置文档缓存失败: %d, error: %s", docID, err)
			// 缓存设置失败不影响返回结果
		} else {
			logc.Infof(ctx, "文档 %d 缓存设置成功", docID)
		}
	}

	return &doc, nil
}

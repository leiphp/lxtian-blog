package utils

import (
	"context"
	"fmt"
	"time"

	redisutil "lxtian-blog/common/pkg/redis"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

// ViewCountUtil 文章浏览次数工具
type ViewCountUtil struct {
	DB  *gorm.DB
	Rds *redis.Redis
}

// NewViewCountUtil 创建浏览次数工具实例
func NewViewCountUtil(db *gorm.DB, rds *redis.Redis) *ViewCountUtil {
	return &ViewCountUtil{
		DB:  db,
		Rds: rds,
	}
}

// IncrementArticleView 增加文章浏览次数（同一IP一天只算一次）
func (v *ViewCountUtil) IncrementArticleView(ctx context.Context, articleID uint32, clientIP string) error {
	// 使用统一的Redis Key管理
	redisKey := redisutil.GetArticleViewKeyToday(articleID, clientIP)

	// 检查今天是否已经访问过
	exists, err := v.Rds.ExistsCtx(ctx, redisKey)
	if err != nil {
		logc.Errorf(ctx, "检查Redis key失败: %s", err)
		return err
	}

	// 如果今天已经访问过，直接返回
	if exists {
		logc.Infof(ctx, "IP %s 今天已经访问过文章 %d", clientIP, articleID)
		return nil
	}

	// 设置Redis key，过期时间为明天凌晨
	expireTime := getTomorrowMidnight()
	err = v.Rds.SetexCtx(ctx, redisKey, "1", int(expireTime.Seconds()))
	if err != nil {
		logc.Errorf(ctx, "设置Redis key失败: %s", err)
		return err
	}

	// 更新数据库中的浏览次数
	err = v.updateArticleViewCount(ctx, articleID)
	if err != nil {
		logc.Errorf(ctx, "更新文章浏览次数失败: %s", err)
		return err
	}

	logc.Infof(ctx, "成功记录IP %s 访问文章 %d", clientIP, articleID)
	return nil
}

// updateArticleViewCount 更新数据库中的文章浏览次数
func (v *ViewCountUtil) updateArticleViewCount(ctx context.Context, articleID uint32) error {
	// 使用原子操作更新浏览次数
	err := v.DB.WithContext(ctx).
		Table("txy_article").
		Where("id = ?", articleID).
		Update("view_count", gorm.Expr("view_count + 1")).Error

	if err != nil {
		return fmt.Errorf("更新文章浏览次数失败: %w", err)
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

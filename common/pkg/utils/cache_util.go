package utils

import (
	"context"
	"fmt"
	"lxtian-blog/common/pkg/redis"

	"github.com/zeromicro/go-zero/core/logx"
	redislib "github.com/zeromicro/go-zero/core/stores/redis"
)

// CacheUtil 缓存工具类
type CacheUtil struct {
	Rds *redislib.Redis
}

// NewCacheUtil 创建缓存工具实例
func NewCacheUtil(rds *redislib.Redis) *CacheUtil {
	return &CacheUtil{
		Rds: rds,
	}
}

// DeleteArticleCache 删除文章缓存
// articleID: 文章ID
func (c *CacheUtil) DeleteArticleCache(ctx context.Context, articleID uint64) error {
	key := c.getArticleCacheKey(articleID)

	_, err := c.Rds.DelCtx(ctx, key)
	if err != nil {
		logx.Errorf("删除文章缓存失败: %v", err)
		return err
	}

	logx.Infof("文章 %d 缓存删除成功", articleID)
	return nil
}

// DeleteChapterCache 删除章节缓存
// chapterID: 章节ID
func (c *CacheUtil) DeleteChapterCache(ctx context.Context, chapterID uint64) error {
	key := redis.ReturnRedisKey(redis.ApiWebStringBookChapter, chapterID)

	_, err := c.Rds.DelCtx(ctx, key)
	if err != nil {
		logx.Errorf("删除章节缓存失败: %v", err)
		return err
	}

	logx.Infof("章节 %d 缓存删除成功", chapterID)
	return nil
}

// getArticleCacheKey 获取文章缓存Key
func (c *CacheUtil) getArticleCacheKey(articleID uint64) string {
	return fmt.Sprintf("%sarticle:detail:%d", redis.KeyPrefix, articleID)
}

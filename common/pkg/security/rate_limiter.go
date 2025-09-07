package security

import (
	"context"
	"fmt"
	"time"

	redisutil "lxtian-blog/common/pkg/redis"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RateLimiter 限流器
type RateLimiter struct {
	Rds *redis.Redis
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rds *redis.Redis) *RateLimiter {
	return &RateLimiter{
		Rds: rds,
	}
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	WindowSize  time.Duration // 时间窗口大小
	MaxRequests int           // 最大请求次数
	KeyPrefix   string        // Redis Key前缀
}

// 默认限流配置
var (
	// 通用接口限流：每分钟最多60次请求
	DefaultRateLimit = RateLimitConfig{
		WindowSize:  time.Minute,
		MaxRequests: 60,
		KeyPrefix:   "rate_limit",
	}

	// 文章详情限流：每分钟最多30次请求
	ArticleRateLimit = RateLimitConfig{
		WindowSize:  time.Minute,
		MaxRequests: 30,
		KeyPrefix:   "article_rate",
	}

	// 分类列表限流：每分钟最多20次请求
	CategoryRateLimit = RateLimitConfig{
		WindowSize:  time.Minute,
		MaxRequests: 20,
		KeyPrefix:   "category_rate",
	}
)

// IsAllowed 检查是否允许访问
func (rl *RateLimiter) IsAllowed(ctx context.Context, clientIP, endpoint string, config RateLimitConfig) (bool, error) {
	// 生成限流Key
	key := rl.generateKey(clientIP, endpoint, config)

	// 获取当前计数
	count, err := rl.Rds.IncrCtx(ctx, key)
	if err != nil {
		logc.Errorf(ctx, "限流检查失败: %s", err)
		return false, err
	}

	// 如果是第一次访问，设置过期时间
	if count == 1 {
		err = rl.Rds.ExpireCtx(ctx, key, int(config.WindowSize.Seconds()))
		if err != nil {
			logc.Errorf(ctx, "设置限流过期时间失败: %s", err)
		}
	}

	// 检查是否超过限制
	allowed := count <= int64(config.MaxRequests)

	if !allowed {
		logc.Errorf(ctx, "IP %s 访问 %s 被限流，当前计数: %d，限制: %d",
			clientIP, endpoint, count, config.MaxRequests)
	}

	return allowed, nil
}

// generateKey 生成限流Key
// 格式: blog:security:{key_prefix}:{ip}:{endpoint}:{time}
func (rl *RateLimiter) generateKey(clientIP, endpoint string, config RateLimitConfig) string {
	// 使用时间窗口的起始时间作为Key的一部分，确保时间窗口的一致性
	windowStart := time.Now().Truncate(config.WindowSize)
	timeKey := windowStart.Format("2006-01-02-15-04")

	return fmt.Sprintf("%ssecurity:%s:%s:%s:%s",
		redisutil.KeyPrefix,
		config.KeyPrefix,
		clientIP,
		endpoint,
		timeKey)
}

// GetRemainingRequests 获取剩余请求次数
func (rl *RateLimiter) GetRemainingRequests(ctx context.Context, clientIP, endpoint string, config RateLimitConfig) (int64, error) {
	key := rl.generateKey(clientIP, endpoint, config)
	count, err := rl.Rds.GetCtx(ctx, key)
	if err != nil {
		return 0, err
	}

	if count == "" {
		return int64(config.MaxRequests), nil
	}

	// 解析当前计数
	var currentCount int64
	_, err = fmt.Sscanf(count, "%d", &currentCount)
	if err != nil {
		return 0, err
	}

	remaining := int64(config.MaxRequests) - currentCount
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// ResetRateLimit 重置限流（管理员功能）
func (rl *RateLimiter) ResetRateLimit(ctx context.Context, clientIP, endpoint string, config RateLimitConfig) error {
	key := rl.generateKey(clientIP, endpoint, config)
	_, err := rl.Rds.DelCtx(ctx, key)
	return err
}

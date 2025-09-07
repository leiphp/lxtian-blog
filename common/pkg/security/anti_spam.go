package security

import (
	"context"
	"fmt"
	"strings"
	"time"

	redisutil "lxtian-blog/common/pkg/redis"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// AntiSpam 反刷接口检测器
type AntiSpam struct {
	Rds *redis.Redis
}

// NewAntiSpam 创建反刷检测器
func NewAntiSpam(rds *redis.Redis) *AntiSpam {
	return &AntiSpam{
		Rds: rds,
	}
}

// SpamConfig 反刷配置
type SpamConfig struct {
	MaxRequestsPerMinute int           // 每分钟最大请求数
	MaxRequestsPerHour   int           // 每小时最大请求数
	BlockDuration        time.Duration // 封禁时长
	KeyPrefix            string        // Redis Key前缀
}

// 默认反刷配置
var (
	DefaultSpamConfig = SpamConfig{
		MaxRequestsPerMinute: 100,            // 每分钟最多100次请求
		MaxRequestsPerHour:   1000,           // 每小时最多1000次请求
		BlockDuration:        time.Hour * 24, // 封禁24小时
		KeyPrefix:            "anti_spam",
	}
)

// CheckSpam 检查是否为恶意请求
func (as *AntiSpam) CheckSpam(ctx context.Context, clientIP, userAgent, endpoint string) (bool, error) {
	// 1. 检查是否已被封禁
	if as.isBlocked(ctx, clientIP) {
		logc.Errorf(ctx, "IP %s 已被封禁，拒绝访问", clientIP)
		return true, nil
	}

	// 2. 检查User-Agent是否可疑
	if as.isSuspiciousUserAgent(userAgent) {
		logc.Errorf(ctx, "IP %s 使用可疑User-Agent: %s", clientIP, userAgent)
		as.recordSuspiciousActivity(ctx, clientIP, "suspicious_user_agent", userAgent)
	}

	// 3. 检查请求频率
	if as.isHighFrequency(ctx, clientIP, endpoint) {
		logc.Errorf(ctx, "IP %s 请求频率过高，疑似刷接口", clientIP)
		as.recordSuspiciousActivity(ctx, clientIP, "high_frequency", endpoint)
		return true, nil
	}

	// 4. 记录正常请求
	as.recordRequest(ctx, clientIP, endpoint)

	return false, nil
}

// isBlocked 检查IP是否被封禁
// 格式: blog:security:block:{ip}
func (as *AntiSpam) isBlocked(ctx context.Context, clientIP string) bool {
	blockKey := fmt.Sprintf("%ssecurity:block:%s", redisutil.KeyPrefix, clientIP)
	exists, err := as.Rds.ExistsCtx(ctx, blockKey)
	if err != nil {
		logc.Errorf(ctx, "检查封禁状态失败: %s", err)
		return false
	}
	return exists
}

// isSuspiciousUserAgent 检查User-Agent是否可疑
func (as *AntiSpam) isSuspiciousUserAgent(userAgent string) bool {
	if userAgent == "" {
		return true
	}

	// 可疑的User-Agent模式
	suspiciousPatterns := []string{
		"bot", "crawler", "spider", "scraper",
		"curl", "wget", "python", "java",
		"postman", "insomnia", "httpie",
	}

	userAgentLower := strings.ToLower(userAgent)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(userAgentLower, pattern) {
			return true
		}
	}

	return false
}

// isHighFrequency 检查是否高频请求
// 格式: blog:security:freq:minute:{ip}:{endpoint}:{minute}
// 格式: blog:security:freq:hour:{ip}:{endpoint}:{hour}
func (as *AntiSpam) isHighFrequency(ctx context.Context, clientIP, endpoint string) bool {
	now := time.Now()

	// 检查分钟级频率
	minuteKey := fmt.Sprintf("%ssecurity:freq:minute:%s:%s:%d",
		redisutil.KeyPrefix, clientIP, endpoint, now.Minute())
	minuteCount, err := as.Rds.IncrCtx(ctx, minuteKey)
	if err != nil {
		logc.Errorf(ctx, "检查分钟频率失败: %s", err)
		return false
	}

	if minuteCount == 1 {
		as.Rds.ExpireCtx(ctx, minuteKey, 60) // 60秒过期
	}

	if minuteCount > int64(DefaultSpamConfig.MaxRequestsPerMinute) {
		return true
	}

	// 检查小时级频率
	hourKey := fmt.Sprintf("%ssecurity:freq:hour:%s:%s:%d",
		redisutil.KeyPrefix, clientIP, endpoint, now.Hour())
	hourCount, err := as.Rds.IncrCtx(ctx, hourKey)
	if err != nil {
		logc.Errorf(ctx, "检查小时频率失败: %s", err)
		return false
	}

	if hourCount == 1 {
		as.Rds.ExpireCtx(ctx, hourKey, 3600) // 3600秒过期
	}

	if hourCount > int64(DefaultSpamConfig.MaxRequestsPerHour) {
		return true
	}

	return false
}

// recordRequest 记录正常请求
// 格式: blog:security:request:{ip}:{timestamp}
func (as *AntiSpam) recordRequest(ctx context.Context, clientIP, endpoint string) {
	// 记录请求日志（可选）
	requestKey := fmt.Sprintf("%ssecurity:request:%s:%d",
		redisutil.KeyPrefix, clientIP, time.Now().Unix())

	requestData := map[string]interface{}{
		"ip":       clientIP,
		"endpoint": endpoint,
		"time":     time.Now().Format("2006-01-02 15:04:05"),
	}

	// 这里可以记录到数据库或日志文件
	_ = requestData
	_ = requestKey
}

// recordSuspiciousActivity 记录可疑活动
// 格式: blog:security:suspicious:{ip}:{timestamp}
func (as *AntiSpam) recordSuspiciousActivity(ctx context.Context, clientIP, activityType, details string) {
	activityKey := fmt.Sprintf("%ssecurity:suspicious:%s:%d",
		redisutil.KeyPrefix, clientIP, time.Now().Unix())

	activityData := map[string]interface{}{
		"ip":       clientIP,
		"activity": activityType,
		"details":  details,
		"time":     time.Now().Format("2006-01-02 15:04:05"),
	}

	// 记录可疑活动
	_ = activityData
	_ = activityKey

	// 如果可疑活动过多，考虑封禁IP
	as.checkAndBlockIP(ctx, clientIP)
}

// checkAndBlockIP 检查并封禁IP
// 格式: blog:security:suspicious_count:{ip}:{hour}
func (as *AntiSpam) checkAndBlockIP(ctx context.Context, clientIP string) {
	// 检查最近1小时内的可疑活动次数
	hourAgo := time.Now().Add(-time.Hour)
	suspiciousKey := fmt.Sprintf("%ssecurity:suspicious_count:%s:%d",
		redisutil.KeyPrefix, clientIP, hourAgo.Hour())

	count, err := as.Rds.IncrCtx(ctx, suspiciousKey)
	if err != nil {
		logc.Errorf(ctx, "检查可疑活动次数失败: %s", err)
		return
	}

	if count == 1 {
		as.Rds.ExpireCtx(ctx, suspiciousKey, 3600) // 1小时过期
	}

	// 如果1小时内可疑活动超过10次，封禁IP
	if count > 10 {
		as.blockIP(ctx, clientIP)
	}
}

// blockIP 封禁IP
// 格式: blog:security:block:{ip}
func (as *AntiSpam) blockIP(ctx context.Context, clientIP string) {
	blockKey := fmt.Sprintf("%ssecurity:block:%s", redisutil.KeyPrefix, clientIP)

	err := as.Rds.SetexCtx(ctx, blockKey, "blocked", int(DefaultSpamConfig.BlockDuration.Seconds()))
	if err != nil {
		logc.Errorf(ctx, "封禁IP失败: %s", err)
		return
	}

	logc.Errorf(ctx, "IP %s 已被封禁 %v", clientIP, DefaultSpamConfig.BlockDuration)
}

// UnblockIP 解封IP（管理员功能）
// 格式: blog:security:block:{ip}
func (as *AntiSpam) UnblockIP(ctx context.Context, clientIP string) error {
	blockKey := fmt.Sprintf("%ssecurity:block:%s", redisutil.KeyPrefix, clientIP)
	_, err := as.Rds.DelCtx(ctx, blockKey)
	return err
}

// GetBlockedIPs 获取被封禁的IP列表（管理员功能）
func (as *AntiSpam) GetBlockedIPs(ctx context.Context) ([]string, error) {
	// 这里需要根据实际情况实现
	// 可以通过扫描Redis中的block:*键来获取
	return []string{}, nil
}

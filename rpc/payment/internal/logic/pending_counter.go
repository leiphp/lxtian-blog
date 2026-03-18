package logic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	redisutil "lxtian-blog/common/pkg/redis"
	"strconv"
	"time"
)

// adjustUserDailyPendingCount 在 Redis 中调整用户某天的「待支付订单」计数。
// delta < 0 时执行安全的减 1：仅当 key 存在且当前值 > 0 才进行 DECR，避免产生负数或错误的 key。
func adjustUserDailyPendingCount(ctx context.Context, rds *redis.Redis, userID int64, createdAt time.Time, delta int64) {
	if rds == nil || delta >= 0 {
		// 当前只在支付成功/取消/关闭时使用减 1 逻辑
		return
	}

	dateStr := createdAt.Format("2006-01-02")
	key := fmt.Sprintf("%spayment:pending:count:%d:%s", redisutil.KeyPrefix, userID, dateStr)

	// 读取当前值，避免对不存在的 key 做 DECR（会变成 -1）
	val, err := rds.GetCtx(ctx, key)
	if err != nil || val == "" {
		return
	}

	current, err := strconv.ParseInt(val, 10, 64)
	if err != nil || current <= 0 {
		return
	}

	_, _ = rds.DecrCtx(ctx, key)
}


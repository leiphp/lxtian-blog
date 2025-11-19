package user_repo

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/model"
	redisutil "lxtian-blog/common/pkg/redis"
	"lxtian-blog/common/repository"
	"time"

	redisstore "github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

// UserMembershipRepository 用户会员表仓储接口
type UserMembershipRepository interface {
	repository.BaseRepository[model.LxtUserMembership]

	// 会员特有方法
	GetActiveMembershipByUserId(ctx context.Context, userID int64) (map[string]interface{}, error)
	GetByUserId(ctx context.Context, userID int64) (*model.LxtUserMembership, error)
	UpdateIsActive(ctx context.Context, userID int64, isActive int32) error
}

// userMembershipRepository 用户会员表仓储实现
type userMembershipRepository struct {
	*repository.TransactionalBaseRepository[model.LxtUserMembership]
	Rds *redisstore.Redis
}

// NewUserMembershipRepository 创建UserMembership仓储
func NewUserMembershipRepository(db *gorm.DB, rds *redisstore.Redis) UserMembershipRepository {
	return &userMembershipRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.LxtUserMembership](db),
		Rds:                         rds,
	}
}

// GetByUserId 根据用户ID获取会员信息
func (r *userMembershipRepository) GetByUserId(ctx context.Context, userID int64) (*model.LxtUserMembership, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"user_id": userID,
	})
}

// UpdateIsActive 更新会员激活状态
func (r *userMembershipRepository) UpdateIsActive(ctx context.Context, userID int64, isActive int32) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"user_id": userID},
		map[string]interface{}{"is_active": isActive},
	)
}

// GetActiveMembershipByUserId 获取用户会员信息并检查是否过期
// 返回会员信息的 map，包括：end_time, level, is_valid, is_active, start_time, total_days, type_id
func (r *userMembershipRepository) GetActiveMembershipByUserId(ctx context.Context, userID int64) (map[string]interface{}, error) {
	// 1. 优先从 Redis 获取
	if r.Rds != nil {
		cacheKey := redisutil.ReturnRedisKey(redisutil.UserMemberShipString, userID)
		cacheData, err := r.Rds.GetCtx(ctx, cacheKey)
		if err != nil && err != redisstore.Nil {
			// Redis 异常时直接返回错误，避免掩盖问题
			return nil, err
		}
		if cacheData != "" {
			var membershipData map[string]interface{}
			if err := json.Unmarshal([]byte(cacheData), &membershipData); err == nil {
				// 缓存命中并解析成功，直接返回
				return membershipData, nil
			}
			// 解析失败则继续走 DB 逻辑，不影响后续
		}
	}

	// 2. Redis 未命中，从 DB 查询
	var membership model.LxtUserMembership
	err := r.GetDB(ctx).Where("user_id = ?", userID).First(&membership).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 没有会员信息，返回 nil
		}
		return nil, err
	}

	now := time.Now()
	// 检查会员是否过期：end_time < 当前时间 && is_active = 1
	if membership.EndTime.Before(now) && membership.IsActive == 1 {
		// 会员已过期，更新 is_active = 0
		err = r.GetDB(ctx).Model(&membership).Update("is_active", 0).Error
		if err != nil {
			return nil, err
		}
		membership.IsActive = 0
	}

	// 判断会员是否有效：end_time > 当前时间 && is_active = 1
	isValid := membership.EndTime.After(now) && membership.IsActive == 1

	// 转换为 map
	membershipData := map[string]interface{}{
		"end_time":     membership.EndTime.Format("2006-01-02 15:04:05"),
		"level":        membership.Level,
		"is_valid":     isValid,
		"is_active":    membership.IsActive,
		"start_time":   membership.StartTime.Format("2006-01-02 15:04:05"),
		"total_months": membership.TotalMonths,
		"type_id":      membership.MembershipTypeID,
	}

	// 3. 将结果写入 Redis，方便下次直接读取
	if r.Rds != nil {
		cacheKey := redisutil.ReturnRedisKey(redisutil.UserMemberShipString, userID)
		if data, err := json.Marshal(membershipData); err == nil {
			// 这里简单设置为 1 小时过期，可根据业务调整为 membership 剩余有效时间等
			_ = r.Rds.SetexCtx(ctx, cacheKey, string(data), 3600)
		}
	}

	return membershipData, nil
}

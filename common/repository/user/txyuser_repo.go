package user

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository"
	"time"

	"gorm.io/gorm"
)

// TxyUserRepository TxyUser表仓储接口
type TxyUserRepository interface {
	repository.BaseRepository[mysql.TxyUser]

	// 用户特有方法
	GetByUid(ctx context.Context, uid uint64) (*mysql.TxyUser, error)
	GetByOpenid(ctx context.Context, openid string) (*mysql.TxyUser, error)
	GetByUnionid(ctx context.Context, unionid string) (*mysql.TxyUser, error)
	GetByType(ctx context.Context, userType uint64, page, pageSize int) ([]*mysql.TxyUser, int64, error)
	GetUsersByLastLoginTime(ctx context.Context, startTime, endTime uint64, page, pageSize int) ([]*mysql.TxyUser, int64, error)

	// 更新方法
	UpdateLastLogin(ctx context.Context, uid uint64, loginTime uint64, loginIp string) error
	UpdateLoginTimes(ctx context.Context, uid uint64) error
	UpdateAccessToken(ctx context.Context, uid uint64, accessToken string) error

	// 统计方法
	GetCountByType(ctx context.Context, userType uint64) (int64, error)
	GetActiveUserCount(ctx context.Context, days int) (int64, error)
	GetTotalLoginCount(ctx context.Context) (int64, error)

	// 批量操作
	BatchUpdateLastLogin(ctx context.Context, uids []uint64, loginTime uint64, loginIp string) error
	GetInactiveUsers(ctx context.Context, days int, limit int) ([]*mysql.TxyUser, error)
}

// txyUserRepository TxyUser表仓储实现
type txyUserRepository struct {
	*repository.TransactionalBaseRepository[mysql.TxyUser]
}

// NewTxyUserRepository 创建TxyUser仓储
func NewTxyUserRepository(db *gorm.DB) TxyUserRepository {
	return &txyUserRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[mysql.TxyUser](db),
	}
}

// GetByUid 根据Uid获取用户
func (r *txyUserRepository) GetByUid(ctx context.Context, uid uint64) (*mysql.TxyUser, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"uid": uid,
	})
}

// GetByOpenid 根据Openid获取用户
func (r *txyUserRepository) GetByOpenid(ctx context.Context, openid string) (*mysql.TxyUser, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"openid": openid,
	})
}

// GetByUnionid 根据Unionid获取用户
func (r *txyUserRepository) GetByUnionid(ctx context.Context, unionid string) (*mysql.TxyUser, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"unionid": unionid,
	})
}

// GetByType 根据用户类型获取用户列表
func (r *txyUserRepository) GetByType(ctx context.Context, userType uint64, page, pageSize int) ([]*mysql.TxyUser, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"type": userType,
	}, page, pageSize)
}

// GetUsersByLastLoginTime 根据最后登录时间获取用户列表
func (r *txyUserRepository) GetUsersByLastLoginTime(ctx context.Context, startTime, endTime uint64, page, pageSize int) ([]*mysql.TxyUser, int64, error) {
	db := r.GetDB(ctx)
	var users []*mysql.TxyUser
	var total int64

	query := db.Where("last_login_time BETWEEN ? AND ?", startTime, endTime)

	// 获取总数
	if err := query.Model(&mysql.TxyUser{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Order("last_login_time DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateLastLogin 更新最后登录信息
func (r *txyUserRepository) UpdateLastLogin(ctx context.Context, uid uint64, loginTime uint64, loginIp string) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"uid": uid},
		map[string]interface{}{
			"last_login_time": loginTime,
			"last_login_ip":   loginIp,
		},
	)
}

// UpdateLoginTimes 更新登录次数
func (r *txyUserRepository) UpdateLoginTimes(ctx context.Context, uid uint64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyUser{}).
		Where("uid = ?", uid).
		Update("login_times", gorm.Expr("login_times + ?", 1)).Error
}

// UpdateAccessToken 更新AccessToken
func (r *txyUserRepository) UpdateAccessToken(ctx context.Context, uid uint64, accessToken string) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"uid": uid},
		map[string]interface{}{"access_token": accessToken},
	)
}

// GetCountByType 根据用户类型统计数量
func (r *txyUserRepository) GetCountByType(ctx context.Context, userType uint64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"type": userType,
	})
}

// GetActiveUserCount 获取活跃用户数量
func (r *txyUserRepository) GetActiveUserCount(ctx context.Context, days int) (int64, error) {
	db := r.GetDB(ctx)
	var count int64

	// 计算时间戳
	cutoffTime := uint64(0)
	if days > 0 {
		cutoffTime = uint64(time.Now().Unix() - int64(days*24*3600))
	}

	err := db.Model(&mysql.TxyUser{}).
		Where("last_login_time > ?", cutoffTime).
		Count(&count).Error

	return count, err
}

// GetTotalLoginCount 获取总登录次数
func (r *txyUserRepository) GetTotalLoginCount(ctx context.Context) (int64, error) {
	db := r.GetDB(ctx)
	var total int64

	err := db.Model(&mysql.TxyUser{}).
		Select("COALESCE(SUM(login_times), 0)").
		Scan(&total).Error

	return total, err
}

// BatchUpdateLastLogin 批量更新最后登录信息
func (r *txyUserRepository) BatchUpdateLastLogin(ctx context.Context, uids []uint64, loginTime uint64, loginIp string) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyUser{}).
		Where("uid IN ?", uids).
		Updates(map[string]interface{}{
			"last_login_time": loginTime,
			"last_login_ip":   loginIp,
		}).Error
}

// GetInactiveUsers 获取不活跃用户
func (r *txyUserRepository) GetInactiveUsers(ctx context.Context, days int, limit int) ([]*mysql.TxyUser, error) {
	db := r.GetDB(ctx)
	var users []*mysql.TxyUser

	// 计算时间戳
	cutoffTime := uint64(time.Now().Unix() - int64(days*24*3600))

	err := db.Where("last_login_time < ? OR last_login_time = 0", cutoffTime).
		Order("last_login_time ASC").
		Limit(limit).
		Find(&users).Error

	return users, err
}

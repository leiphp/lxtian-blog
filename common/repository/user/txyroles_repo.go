package user

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository"

	"gorm.io/gorm"
)

// TxyRolesRepository TxyRoles表仓储接口
type TxyRolesRepository interface {
	repository.BaseRepository[mysql.TxyRoles]

	// 角色特有方法
	GetByKey(ctx context.Context, key string) (*mysql.TxyRoles, error)
	GetByStatus(ctx context.Context, status int64) ([]*mysql.TxyRoles, error)
	GetActiveRoles(ctx context.Context) ([]*mysql.TxyRoles, error)

	// 更新方法
	UpdateStatus(ctx context.Context, roleId uint64, status int64) error
	UpdateDescription(ctx context.Context, roleId uint64, description string) error

	// 统计方法
	GetCountByStatus(ctx context.Context, status int64) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, roleIds []uint64, status int64) error
}

// txyRolesRepository TxyRoles表仓储实现
type txyRolesRepository struct {
	*repository.TransactionalBaseRepository[mysql.TxyRoles]
}

// NewTxyRolesRepository 创建TxyRoles仓储
func NewTxyRolesRepository(db *gorm.DB) TxyRolesRepository {
	return &txyRolesRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[mysql.TxyRoles](db),
	}
}

// GetByKey 根据Key获取角色
func (r *txyRolesRepository) GetByKey(ctx context.Context, key string) (*mysql.TxyRoles, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"key": key,
	})
}

// GetByStatus 根据状态获取角色列表
func (r *txyRolesRepository) GetByStatus(ctx context.Context, status int64) ([]*mysql.TxyRoles, error) {
	roles, _, err := r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, 0, 0) // 不分页
	return roles, err
}

// GetActiveRoles 获取启用的角色列表
func (r *txyRolesRepository) GetActiveRoles(ctx context.Context) ([]*mysql.TxyRoles, error) {
	return r.GetByStatus(ctx, 1) // 1表示启用状态
}

// UpdateStatus 更新角色状态
func (r *txyRolesRepository) UpdateStatus(ctx context.Context, roleId uint64, status int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": roleId},
		map[string]interface{}{"status": status},
	)
}

// UpdateDescription 更新角色描述
func (r *txyRolesRepository) UpdateDescription(ctx context.Context, roleId uint64, description string) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": roleId},
		map[string]interface{}{"description": description},
	)
}

// GetCountByStatus 根据状态统计角色数量
func (r *txyRolesRepository) GetCountByStatus(ctx context.Context, status int64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"status": status,
	})
}

// BatchUpdateStatus 批量更新角色状态
func (r *txyRolesRepository) BatchUpdateStatus(ctx context.Context, roleIds []uint64, status int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyRoles{}).
		Where("id IN ?", roleIds).
		Update("status", status).Error
}

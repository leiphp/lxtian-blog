package user_repo

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository"

	"gorm.io/gorm"
)

// TxyPermissionsRepository TxyPermissions表仓储接口
type TxyPermissionsRepository interface {
	repository.BaseRepository[mysql.TxyPermissions]

	// 权限特有方法
	GetByPermissionCode(ctx context.Context, permissionCode string) (*mysql.TxyPermissions, error)
	GetByModule(ctx context.Context, module string) ([]*mysql.TxyPermissions, error)
	GetByStatus(ctx context.Context, status int64) ([]*mysql.TxyPermissions, error)
	GetActivePermissions(ctx context.Context) ([]*mysql.TxyPermissions, error)

	// 更新方法
	UpdateStatus(ctx context.Context, permId uint64, status int64) error
	UpdateDescription(ctx context.Context, permId uint64, description string) error

	// 统计方法
	GetCountByStatus(ctx context.Context, status int64) (int64, error)
	GetCountByModule(ctx context.Context, module string) (int64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, permIds []uint64, status int64) error
}

// txyPermissionsRepository TxyPermissions表仓储实现
type txyPermissionsRepository struct {
	*repository.TransactionalBaseRepository[mysql.TxyPermissions]
}

// NewTxyPermissionsRepository 创建TxyPermissions仓储
func NewTxyPermissionsRepository(db *gorm.DB) TxyPermissionsRepository {
	return &txyPermissionsRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[mysql.TxyPermissions](db),
	}
}

// GetByPermissionCode 根据权限码获取权限
func (r *txyPermissionsRepository) GetByPermissionCode(ctx context.Context, permissionCode string) (*mysql.TxyPermissions, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"permission_code": permissionCode,
	})
}

// GetByModule 根据模块获取权限列表
func (r *txyPermissionsRepository) GetByModule(ctx context.Context, module string) ([]*mysql.TxyPermissions, error) {
	permissions, _, err := r.GetList(ctx, map[string]interface{}{
		"module": module,
	}, 0, 0, "", "") // 不分页
	return permissions, err
}

// GetByStatus 根据状态获取权限列表
func (r *txyPermissionsRepository) GetByStatus(ctx context.Context, status int64) ([]*mysql.TxyPermissions, error) {
	permissions, _, err := r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, 0, 0, "", "") // 不分页
	return permissions, err
}

// GetActivePermissions 获取启用的权限列表
func (r *txyPermissionsRepository) GetActivePermissions(ctx context.Context) ([]*mysql.TxyPermissions, error) {
	return r.GetByStatus(ctx, 1) // 1表示启用状态
}

// UpdateStatus 更新权限状态
func (r *txyPermissionsRepository) UpdateStatus(ctx context.Context, permId uint64, status int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": permId},
		map[string]interface{}{"status": status},
	)
}

// UpdateDescription 更新权限描述
func (r *txyPermissionsRepository) UpdateDescription(ctx context.Context, permId uint64, description string) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": permId},
		map[string]interface{}{"description": description},
	)
}

// GetCountByStatus 根据状态统计权限数量
func (r *txyPermissionsRepository) GetCountByStatus(ctx context.Context, status int64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"status": status,
	})
}

// GetCountByModule 根据模块统计权限数量
func (r *txyPermissionsRepository) GetCountByModule(ctx context.Context, module string) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"module": module,
	})
}

// BatchUpdateStatus 批量更新权限状态
func (r *txyPermissionsRepository) BatchUpdateStatus(ctx context.Context, permIds []uint64, status int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyPermissions{}).
		Where("id IN ?", permIds).
		Update("status", status).Error
}

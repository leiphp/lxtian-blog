package user_repo

import (
	"context"
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository"

	"gorm.io/gorm"
)

// MembershipPermissionRepository 会员权限仓储接口
type MembershipPermissionRepository interface {
	repository.BaseRepository[model.LxtUserMembershipPermission]

	// 权限特有方法
	FindByMembershipTypeId(ctx context.Context, membershipTypeId uint64) ([]*model.LxtUserMembershipPermission, error)
	GetPermissionKeysByTypeId(ctx context.Context, membershipTypeId uint64) ([]string, error)
	BatchCreateByTypeId(ctx context.Context, membershipTypeId int64, permissions []*model.LxtUserMembershipPermission) error
}

// membershipPermissionRepository 会员权限仓储实现
type membershipPermissionRepository struct {
	*repository.TransactionalBaseRepository[model.LxtUserMembershipPermission]
}

// NewMembershipPermissionRepository 创建会员权限仓储
func NewMembershipPermissionRepository(db *gorm.DB) MembershipPermissionRepository {
	return &membershipPermissionRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.LxtUserMembershipPermission](db),
	}
}

// FindByMembershipTypeId 根据会员类型ID查询权限
func (r *membershipPermissionRepository) FindByMembershipTypeId(ctx context.Context, membershipTypeId uint64) ([]*model.LxtUserMembershipPermission, error) {
	var entities []*model.LxtUserMembershipPermission
	db := r.GetDB(ctx)
	err := db.
		Where("membership_type_id = ? AND deleted_at IS NULL", membershipTypeId).
		Find(&entities).Error

	return entities, err
}

// GetPermissionKeysByTypeId 根据会员类型ID获取权限键列表
func (r *membershipPermissionRepository) GetPermissionKeysByTypeId(ctx context.Context, membershipTypeId uint64) ([]string, error) {
	permissions, err := r.FindByMembershipTypeId(ctx, membershipTypeId)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(permissions))
	for _, p := range permissions {
		keys = append(keys, p.PermissionKey)
	}
	return keys, nil
}

// BatchCreateByTypeId 批量创建权限
func (r *membershipPermissionRepository) BatchCreateByTypeId(ctx context.Context, membershipTypeId int64, permissions []*model.LxtUserMembershipPermission) error {
	for _, p := range permissions {
		p.MembershipTypeID = membershipTypeId
	}
	db := r.GetDB(ctx)
	return db.CreateInBatches(permissions, 100).Error
}

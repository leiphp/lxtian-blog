package user

import (
	"context"
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository"

	"gorm.io/gorm"
)

// MembershipTypeRepository 会员类型仓储接口
type MembershipTypeRepository interface {
	repository.BaseRepository[model.LxtUserMembershipTypes]

	// 会员类型特有方法
	FindAllActive(ctx context.Context) ([]*model.LxtUserMembershipTypes, error)
	GetByKey(ctx context.Context, key string) (*model.LxtUserMembershipTypes, error)
	//GetList(ctx context.Context, page, pageSize int, status int64) ([]*model.LxtUserMembershipTypes, int64, error)
}

// membershipTypeRepository 会员类型仓储实现
type membershipTypeRepository struct {
	*repository.TransactionalBaseRepository[model.LxtUserMembershipTypes]
}

// NewMembershipTypeRepository 创建会员类型仓储
func NewMembershipTypeRepository(db *gorm.DB) MembershipTypeRepository {
	return &membershipTypeRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.LxtUserMembershipTypes](db),
	}
}

// FindAllActive 查询所有启用的会员类型
func (r *membershipTypeRepository) FindAllActive(ctx context.Context) ([]*model.LxtUserMembershipTypes, error) {
	var entities []*model.LxtUserMembershipTypes
	db := r.GetDB(ctx)
	err := db.WithContext(ctx).
		Where("status = ? AND deleted_at IS NULL", 1).
		Order("id ASC").
		Find(&entities).Error

	return entities, err
}

// GetByKey 根据key获取会员类型
func (r *membershipTypeRepository) GetByKey(ctx context.Context, key string) (*model.LxtUserMembershipTypes, error) {
	var entity model.LxtUserMembershipTypes
	db := r.GetDB(ctx)
	err := db.WithContext(ctx).
		Where("`key` = ? AND deleted_at IS NULL", key).
		First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// GetList 获取会员类型列表（支持分页和状态筛选）
//func (r *membershipTypeRepository) GetList(ctx context.Context, page, pageSize int, status int64) ([]*model.LxtUserMembershipTypes, int64, error) {
//	var entities []*model.LxtUserMembershipTypes
//	var total int64
//
//	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")
//
//	// 状态筛选
//	if status >= 0 {
//		query = query.Where("status = ?", status)
//	}
//
//	// 获取总数
//	if err := query.Model(&model.LxtUserMembershipTypes{}).Count(&total).Error; err != nil {
//		return nil, 0, err
//	}
//
//	// 分页查询
//	if page > 0 && pageSize > 0 {
//		offset := (page - 1) * pageSize
//		query = query.Offset(offset).Limit(pageSize)
//	}
//
//	err := query.Order("id DESC").Find(&entities).Error
//	return entities, total, err
//}

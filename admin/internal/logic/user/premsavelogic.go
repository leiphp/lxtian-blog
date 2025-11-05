package user

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PremSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 权限保存
func NewPremSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PremSaveLogic {
	return &PremSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PremSaveLogic) PremSave(req *types.PremSaveReq) (resp *types.PremSaveResp, err error) {
	db := l.svcCtx.DB

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 删除旧权限
	if err = tx.Delete(&mysql.TxyRolePermissions{}, req.RoleId).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. 批量插入新权限
	var rolePermissions []mysql.TxyRolePermissions
	for _, permId := range req.PermIds {
		rolePermissions = append(rolePermissions, mysql.TxyRolePermissions{
			RoleId: uint64(req.RoleId),
			PermId: uint64(permId),
		})
	}
	if err := tx.Create(&rolePermissions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 3. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	resp = new(types.PremSaveResp)
	resp.Data = true
	return
}

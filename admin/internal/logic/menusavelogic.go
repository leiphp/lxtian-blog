package logic

import (
	"context"
	"database/sql"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenuSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 菜单保存
func NewMenuSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenuSaveLogic {
	return &MenuSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenuSaveLogic) MenuSave(req *types.MenuSaveReq) (resp *types.MenuSaveResp, err error) {
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

	// 1.判断是插入还是更新
	menu := mysql.TxyMenu{
		Title:   req.Title,
		Pid:     req.Pid,
		Index:   req.Index,
		Icon:    req.Icon,
		Permiss: req.Permiss,
		Sort:    req.Sort,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	if req.Id == 0 {
		menu.CreatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		if err = tx.Create(&menu).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		perm := mysql.TxyPermissions{
			Name:     req.Title,
			Type:     "menu",
			ParentId: req.Pid,
			MenuId:   int64(menu.Id),
			CreatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
		}
		// 插入权限表id
		if err = tx.Create(&perm).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		menu.Id = uint64(req.Id)
		if err = tx.Model(&menu).Updates(menu).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	// 3. 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	resp = new(types.MenuSaveResp)
	resp.Data = true
	return
}

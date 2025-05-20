package logic

import (
	"context"
	"errors"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"

	"github.com/zeromicro/go-zero/core/logx"
)

type MenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 菜单管理
func NewMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MenusLogic {
	return &MenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MenusLogic) Menus(req *types.MenusReq) (resp *types.MenusResp, err error) {
	if req.Perm != "menu" && req.Perm != "" {
		return nil, errors.New("invalid permission")
	}
	resp = new(types.MenusResp)
	db := l.svcCtx.DB.Table("txy_menu").Order("txy_menu.sort ASC")
	if req.Perm != "" {
		var results []MenuWithPerm
		err = db.Select("txy_menu.id, txy_menu.title, txy_menu.pid, txy_menu.`index`, txy_menu.icon, txy_menu.permiss, txy_menu.sort, txy_permissions.id as perm_id").
			Joins("LEFT JOIN txy_permissions ON txy_permissions.menu_id = txy_menu.id AND txy_permissions.type = ?", "menu").
			Group("txy_menu.id").
			Scan(&results).Error
		if err != nil {
			return nil, err
		}

		for _, item := range results {
			resp.Data = append(resp.Data, map[string]interface{}{
				"id":      item.Id,
				"title":   item.Title,
				"pid":     item.Pid,
				"index":   item.Index,
				"icon":    item.Icon,
				"permiss": item.Permiss,
				"sort":    item.Sort,
				"perm_id": item.PermId,
			})
		}
	} else {
		var results []map[string]interface{}
		err = db.Select("id, title, pid, `index`, icon, permiss, sort").
			Order("sort ASC").
			Find(&results).Error
		if err != nil {
			return nil, err
		}
		resp.Data = results
	}
	return
}

type MenuWithPerm struct {
	mysql.TxyMenu
	PermId int64 `json:"perm_id"`
}

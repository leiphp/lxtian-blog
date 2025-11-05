package user

import (
	"context"
	"fmt"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type RolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 角色管理
func NewRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RolesLogic {
	return &RolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RolesLogic) Roles(req *types.RolesReq) (resp *types.RolesResp, err error) {
	// 基础查询构建（包含JOIN和公共WHERE条件）
	baseDB := l.svcCtx.DB.Model(&mysql.TxyRoles{})
	order := "txy_roles.id asc"
	// 填充WHERE条件
	if req.Keywords != "" {
		baseDB = baseDB.Where("txy_roles.name like ?", "%"+req.Keywords+"%")
	}
	// 计算总数（使用基础查询，无分页/排序）
	var total int64
	if err := baseDB.Count(&total).Error; err != nil {
		return nil, err
	}

	// 处理分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize

	// 1. 获取角色列表
	var results []map[string]interface{}
	err = baseDB.Select("txy_roles.id,txy_roles.name,txy_roles.description,txy_roles.status,txy_roles.key,txy_roles.created_at,txy_roles.updated_at").
		Limit(req.PageSize).
		Offset(offset).
		Order(order).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 2. 提取角色ID列表
	roleIDs := make([]int64, 0)
	for _, r := range results {
		val := r["id"]
		switch v := val.(type) {
		case int64:
			roleIDs = append(roleIDs, v)
		case uint64:
			roleIDs = append(roleIDs, int64(v))
		case int:
			roleIDs = append(roleIDs, int64(v))
		case float64:
			roleIDs = append(roleIDs, int64(v))
		case string:
			if idParsed, err := strconv.ParseInt(v, 10, 64); err == nil {
				roleIDs = append(roleIDs, idParsed)
			}
		default:
			fmt.Printf("unrecognized id type: %T => %#v\n", val, val)
		}
	}

	// 3. 查询对应角色的 menu_id
	type RolePermission struct {
		RoleID int64
		PermID int64
	}
	var perms []RolePermission
	err = l.svcCtx.DB.
		Table("txy_role_permissions as rp").
		Select("rp.role_id, rp.perm_id").
		Where("rp.role_id IN ?", roleIDs).Debug().
		Scan(&perms).Error
	if err != nil {
		return nil, err
	}

	// 4. 构建角色ID => []menu_id 映射
	permMap := make(map[int64][]int64)
	for _, p := range perms {
		permMap[p.RoleID] = append(permMap[p.RoleID], p.PermID)
	}

	// 5. 把权限字段塞入每个角色记录
	for _, r := range results {
		if uid, ok := r["id"].(uint64); ok {
			id := int64(uid)
			if perms, ok := permMap[id]; ok {
				r["permiss"] = perms
			} else {
				r["permiss"] = []any{} // 保证是空数组，而不是 null
			}
		}
	}

	// 转换 []byte -> string（特别是中文字段）
	utils.ConvertByteFieldsToString(results)

	resp = new(types.RolesResp)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = total
	resp.List = results
	return
}

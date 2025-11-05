package user

import (
	"context"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type UsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户管理
func NewUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UsersLogic {
	return &UsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UsersLogic) Users(req *types.UsersReq) (resp *types.UsersResp, err error) {
	// 基础查询构建（包含JOIN和公共WHERE条件）
	baseDB := l.svcCtx.DB.Model(&mysql.TxyUser{}).
		Joins("left join txy_user_roles as ur on ur.user_id = txy_user.id").
		Joins("left join txy_roles as r on r.id = ur.role_id")
	order := "txy_user.id desc"
	// 填充WHERE条件
	if req.Role != "" {
		baseDB = baseDB.Where("r.`key` = ?", req.Role)
	}
	if req.Keywords != "" {
		baseDB = baseDB.Where("txy_user.username like ?", "%"+req.Keywords+"%")
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
	var users []map[string]interface{}
	err = baseDB.Select("txy_user.id,txy_user.username,txy_user.password,txy_user.nickname,txy_user.created_at,txy_user.email,txy_user.head_img,r.name role_name").
		Limit(req.PageSize).
		Offset(offset).
		Order(order).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	// 转换 []byte -> string（特别是中文字段）
	utils.ConvertByteFieldsToString(users)

	resp = new(types.UsersResp)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.Total = total
	resp.List = users
	return
}

package userlogic

import (
	"context"
	"database/sql"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMembershipListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMembershipListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMembershipListLogic {
	return &GetMembershipListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会员列表
func (l *GetMembershipListLogic) GetMembershipList(in *user.GetMembershipListReq) (*user.GetMembershipListResp, error) {
	// 查询所有启用的会员类型
	type MembershipTypeRow struct {
		Id          int64
		Name        string
		Key         string
		Days        int64
		Description sql.NullString
		Status      int64
		CreatedAt   string
		UpdatedAt   string
	}

	var membershipTypes []MembershipTypeRow

	// 使用GORM查询启用的会员类型
	err := l.svcCtx.DB.Table("lxt_membership_types").
		Where("status = ? AND deleted_at IS NULL", 1).
		Order("id ASC").
		Find(&membershipTypes).
		Error

	if err != nil {
		l.Errorf("查询会员类型失败: %v", err)
		return nil, err
	}

	// 查询会员权限
	type PermissionRow struct {
		MembershipTypeId int64
		PermissionKey    string
	}

	var permissions []PermissionRow
	err = l.svcCtx.DB.Table("lxt_membership_permissions").
		Select("membership_type_id, permission_key").
		Find(&permissions).
		Error

	if err != nil {
		l.Errorf("查询会员权限失败: %v", err)
	}

	// 按会员类型ID分组权限
	permissionMap := make(map[int64][]string)
	for _, p := range permissions {
		permissionMap[p.MembershipTypeId] = append(permissionMap[p.MembershipTypeId], p.PermissionKey)
	}

	// 转换为响应格式
	respList := make([]*user.MembershipType, 0, len(membershipTypes))
	for _, mt := range membershipTypes {
		// 获取描述
		description := ""
		if mt.Description.Valid {
			description = mt.Description.String
		}

		// 获取价格信息（根据天数计算）
		price := l.calculatePrice(mt.Days)
		originalPrice := price * 1.2 // 原价是现价的1.2倍
		discount := 1.0 - (price / originalPrice)
		if discount < 0 {
			discount = 0
		}

		// 判断是否推荐（默认年度会员推荐）
		popular := false
		if mt.Key == "yearly" {
			popular = true
		}

		// 获取权限列表
		permissionList := permissionMap[mt.Id]
		if permissionList == nil {
			permissionList = []string{}
		}

		respList = append(respList, &user.MembershipType{
			Id:            int64(mt.Id),
			Name:          mt.Name,
			Price:         price,
			OriginalPrice: originalPrice,
			Discount:      discount,
			Period:        l.getPeriod(mt.Days),
			Popular:       popular,
			Permissions:   permissionList,
			Description:   description,
		})
	}

	return &user.GetMembershipListResp{
		List: respList,
	}, nil
}

// 根据天数计算价格
func (l *GetMembershipListLogic) calculatePrice(days int64) float64 {
	// 基准：月度会员（30天）= 30元
	// 季度会员（92天）= 89元
	// 年度会员（365天）= 299元
	switch {
	case days <= 30:
		return 30.0
	case days <= 92:
		return 89.0
	case days <= 365:
		return 299.0
	default:
		return 30.0
	}
}

// 获取周期单位
func (l *GetMembershipListLogic) getPeriod(days int64) string {
	switch {
	case days <= 30:
		return "月"
	case days <= 92:
		return "季度"
	case days <= 365:
		return "年"
	default:
		return "月"
	}
}

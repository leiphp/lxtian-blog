package payment

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/repository/user_repo"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MembershipListLogic struct {
	logx.Logger
	ctx                       context.Context
	svcCtx                    *svc.ServiceContext
	membershipTypesRepository user_repo.MembershipTypeRepository
}

// 会员套餐
func NewMembershipListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MembershipListLogic {
	return &MembershipListLogic{
		Logger:                    logx.WithContext(ctx),
		ctx:                       ctx,
		svcCtx:                    svcCtx,
		membershipTypesRepository: user_repo.NewMembershipTypeRepository(svcCtx.DB),
	}
}

func (l *MembershipListLogic) MembershipList() (resp *types.MembershipListResp, err error) {
	// 使用Repository查询所有启用的会员类型
	membershipTypes, err := l.membershipTypesRepository.FindAllActive(l.ctx)
	if err != nil {
		l.Errorf("查询会员类型失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	respList := make([]*types.MembershipType, 0, len(membershipTypes))
	for _, mt := range membershipTypes {
		// 解析 JSON 权限字段
		var permissionsJSON string
		if mt.Permissions != nil {
			permissionsJSON = *mt.Permissions
		}
		permissions := l.parsePermissions(permissionsJSON)

		// 处理描述字段
		description := ""
		if mt.Description != nil {
			description = *mt.Description
		}

		respList = append(respList, &types.MembershipType{
			Id:            mt.ID,
			Name:          mt.Name,
			Price:         mt.Price,
			OriginalPrice: mt.OriginalPrice,
			Discount:      mt.Discount,
			Period:        mt.Period,
			Popular:       mt.Popular == 1,
			Permissions:   permissions,
			Description:   description,
		})
	}

	return &types.MembershipListResp{
		List: respList,
	}, nil

}

// parsePermissions 解析权限 JSON 字符串
func (l *MembershipListLogic) parsePermissions(permissionsJSON string) []string {
	if permissionsJSON == "" {
		return []string{}
	}

	var permissions []string
	if err := json.Unmarshal([]byte(permissionsJSON), &permissions); err != nil {
		l.Errorf("解析权限JSON失败: %v, JSON内容: %s", err, permissionsJSON)
		return []string{}
	}

	return permissions
}

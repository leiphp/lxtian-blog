package userlogic

import (
	"context"
	"encoding/json"
	userrepo "lxtian-blog/common/repository/user"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMembershipListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	membershipTypesRepository       userrepo.MembershipTypeRepository
	membershipPermissionsRepository userrepo.MembershipPermissionRepository
}

func NewGetMembershipListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMembershipListLogic {
	return &GetMembershipListLogic{
		ctx:                             ctx,
		svcCtx:                          svcCtx,
		Logger:                          logx.WithContext(ctx),
		membershipTypesRepository:       userrepo.NewMembershipTypeRepository(svcCtx.DB),
		membershipPermissionsRepository: userrepo.NewMembershipPermissionRepository(svcCtx.DB),
	}
}

// 获取会员列表
func (l *GetMembershipListLogic) GetMembershipList(in *user.GetMembershipListReq) (*user.GetMembershipListResp, error) {
	// 使用Repository查询所有启用的会员类型
	membershipTypes, err := l.membershipTypesRepository.FindAllActive(l.ctx)
	if err != nil {
		l.Errorf("查询会员类型失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	respList := make([]*user.MembershipType, 0, len(membershipTypes))
	for _, mt := range membershipTypes {
		// 解析 JSON 权限字段
		var permissionsJSON string
		if mt.Permissions.Valid {
			permissionsJSON = mt.Permissions.String
		}
		permissions := l.parsePermissions(permissionsJSON)

		// 处理描述字段
		description := ""
		if mt.Description.Valid {
			description = mt.Description.String
		}

		respList = append(respList, &user.MembershipType{
			Id:            int64(mt.Id),
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

	return &user.GetMembershipListResp{
		List: respList,
	}, nil
}

// parsePermissions 解析权限 JSON 字符串
func (l *GetMembershipListLogic) parsePermissions(permissionsJSON string) []string {
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

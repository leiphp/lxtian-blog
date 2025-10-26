package user

import (
	"context"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logc"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMembershipListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取会员列表
func NewGetMembershipListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMembershipListLogic {
	return &GetMembershipListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMembershipListLogic) GetMembershipList() (resp *types.GetMembershipListResp, err error) {
	res, err := l.svcCtx.UserRpc.GetMembershipList(l.ctx, &user.GetMembershipListReq{})
	if err != nil {
		logc.Errorf(l.ctx, "GetMembershipList error message: %s", err)
		return nil, err
	}

	// 转换 rpc 响应到 gateway 响应
	resp = &types.GetMembershipListResp{
		List: make([]*types.MembershipType, 0, len(res.List)),
	}

	for _, mt := range res.List {
		gatewayMt := &types.MembershipType{
			Id:            mt.Id,
			Name:          mt.Name,
			Price:         mt.Price,
			OriginalPrice: mt.OriginalPrice,
			Discount:      mt.Discount,
			Period:        mt.Period,
			Popular:       mt.Popular,
			Permissions:   mt.Permissions,
			Description:   mt.Description,
		}
		resp.List = append(resp.List, gatewayMt)
	}

	return resp, nil
}

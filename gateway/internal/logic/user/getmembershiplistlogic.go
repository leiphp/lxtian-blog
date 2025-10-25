package user

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}

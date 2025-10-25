package user

import (
	"context"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpgradeMembershipLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 升级/续费会员
func NewUpgradeMembershipLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpgradeMembershipLogic {
	return &UpgradeMembershipLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpgradeMembershipLogic) UpgradeMembership(req *types.UpgradeMembershipReq) (resp *types.UpgradeMembershipResp, err error) {
	// todo: add your logic here and delete this line

	return
}

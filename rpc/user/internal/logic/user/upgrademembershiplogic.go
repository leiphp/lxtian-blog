package userlogic

import (
	"context"

	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpgradeMembershipLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpgradeMembershipLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpgradeMembershipLogic {
	return &UpgradeMembershipLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpgradeMembershipLogic) UpgradeMembership(in *user.UpgradeMembershipReq) (*user.UpgradeMembershipResp, error) {
	// todo: add your logic here and delete this line

	return &user.UpgradeMembershipResp{}, nil
}

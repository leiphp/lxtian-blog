package payment

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManualRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 手动退款
func NewManualRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManualRefundLogic {
	return &ManualRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManualRefundLogic) ManualRefund(req *types.ManualRefundReq) (resp *types.ManualRefundResp, err error) {
	// todo: add your logic here and delete this line

	return
}

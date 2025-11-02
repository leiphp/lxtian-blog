package logic

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付订单管理
func NewPaymentOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentOrdersLogic {
	return &PaymentOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentOrdersLogic) PaymentOrders(req *types.PaymentOrdersReq) (resp *types.PaymentOrdersResp, err error) {
	// todo: add your logic here and delete this line

	return
}

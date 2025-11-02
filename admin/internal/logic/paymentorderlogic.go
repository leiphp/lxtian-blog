package logic

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付订单详情
func NewPaymentOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentOrderLogic {
	return &PaymentOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentOrderLogic) PaymentOrder(req *types.PaymentOrderReq) (resp *types.PaymentOrderResp, err error) {
	// todo: add your logic here and delete this line

	return
}

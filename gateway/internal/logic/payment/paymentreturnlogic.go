package payment

import (
	"context"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentReturnLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付成功页面跳转
func NewPaymentReturnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentReturnLogic {
	return &PaymentReturnLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentReturnLogic) PaymentReturn(req *types.PaymentNotifyReq) (resp *types.PaymentNotifyResp, err error) {
	// todo: add your logic here and delete this line

	return
}

package payment

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentNotifiesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付通知记录
func NewPaymentNotifiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentNotifiesLogic {
	return &PaymentNotifiesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentNotifiesLogic) PaymentNotifies(req *types.PaymentNotifiesReq) (resp *types.PaymentNotifiesResp, err error) {
	// todo: add your logic here and delete this line

	return
}

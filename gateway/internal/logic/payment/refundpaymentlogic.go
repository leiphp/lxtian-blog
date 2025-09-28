package payment

import (
	"context"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 申请退款
func NewRefundPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundPaymentLogic {
	return &RefundPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundPaymentLogic) RefundPayment(req *types.RefundPaymentReq) (resp *types.RefundPaymentResp, err error) {
	// todo: add your logic here and delete this line

	return
}

package payment

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentRefundsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 退款记录管理
func NewPaymentRefundsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentRefundsLogic {
	return &PaymentRefundsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentRefundsLogic) PaymentRefunds(req *types.PaymentRefundsReq) (resp *types.PaymentRefundsResp, err error) {
	// todo: add your logic here and delete this line

	return
}

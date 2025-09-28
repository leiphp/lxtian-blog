package payment

import (
	"context"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付记录查询
func NewPaymentHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentHistoryLogic {
	return &PaymentHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentHistoryLogic) PaymentHistory(req *types.PaymentHistoryReq) (resp *types.PaymentHistoryResp, err error) {
	// todo: add your logic here and delete this line

	return
}

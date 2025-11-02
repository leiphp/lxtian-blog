package logic

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付统计
func NewPaymentStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentStatsLogic {
	return &PaymentStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentStatsLogic) PaymentStats(req *types.PaymentStatsReq) (resp *types.PaymentStatsResp, err error) {
	// todo: add your logic here and delete this line

	return
}

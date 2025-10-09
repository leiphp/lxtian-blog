package logic

import (
	"context"

	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrdersStatisticsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrdersStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrdersStatisticsLogic {
	return &OrdersStatisticsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 支付订单统计
func (l *OrdersStatisticsLogic) OrdersStatistics(in *payment.OrdersStatisticsReq) (*payment.OrdersStatisticsResp, error) {
	// todo: add your logic here and delete this line

	return &payment.OrdersStatisticsResp{}, nil
}

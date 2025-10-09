package payment

import (
	"context"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrdersStatisticsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付订单统计
func NewOrdersStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrdersStatisticsLogic {
	return &OrdersStatisticsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrdersStatisticsLogic) OrdersStatistics() (resp *types.OrdersStatisticsResp, err error) {
	// todo: add your logic here and delete this line

	return
}

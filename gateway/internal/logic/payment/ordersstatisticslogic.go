package payment

import (
	"context"
	"errors"
	"lxtian-blog/rpc/payment/pb/payment"

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
	//从中间件获取用户信息
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		return nil, errors.New("user_id not found in context")
	}
	res, err := l.svcCtx.PaymentRpc.OrdersStatistics(l.ctx, &payment.OrdersStatisticsReq{
		UserId: int64(userId),
	})
	if err != nil {
		return nil, err
	}

	return &types.OrdersStatisticsResp{
		Total:     res.Total,
		Paid:      res.Paid,
		Pending:   res.Pending,
		PayAmount: float64(res.PayAmount),
	}, nil
}

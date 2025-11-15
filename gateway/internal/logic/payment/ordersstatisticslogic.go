package payment

import (
	"context"
	"errors"
	"lxtian-blog/rpc/payment/pb/payment"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/shopspring/decimal"
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
	// 使用 decimal 库进行精确转换，避免 float32 转 float64 时的精度问题
	// 先将 float32 转换为 float64，然后转换为 decimal，保留2位小数，再转换回 float64
	payAmountDecimal := decimal.NewFromFloat(float64(res.PayAmount)).Round(2)
	payAmountFloat64, _ := payAmountDecimal.Float64()

	return &types.OrdersStatisticsResp{
		Total:       res.Total,
		Paid:        res.Paid,
		Pending:     res.Pending,
		TotalAmount: payAmountFloat64,
	}, nil
}

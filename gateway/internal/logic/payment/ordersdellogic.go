package payment

import (
	"context"
	"errors"
	"lxtian-blog/rpc/payment/pb/payment"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrdersDelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除订单
func NewOrdersDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrdersDelLogic {
	return &OrdersDelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrdersDelLogic) OrdersDel(req *types.OrdersDelReq) (resp *types.OrdersDelResp, err error) {
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		logx.Errorf("OrdersDel user_id not found in context")
		return nil, errors.New("user_id not found in context")
	}
	res, err := l.svcCtx.PaymentRpc.DeletePayment(l.ctx, &payment.DeletePaymentReq{
		UserId:  int64(userId),
		OrderSn: req.OrderSn,
	})
	if err != nil {
		return nil, err
	}

	return &types.OrdersDelResp{
		Message: res.Message,
	}, nil
}

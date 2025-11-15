package payment

import (
	"context"
	"errors"
	"lxtian-blog/rpc/payment/pb/payment"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrdersCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消订单
func NewOrdersCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrdersCancelLogic {
	return &OrdersCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrdersCancelLogic) OrdersCancel(req *types.OrdersCancelReq) (resp *types.OrdersCancelResp, err error) {
	userId, ok := l.ctx.Value("user_id").(int64)
	if !ok {
		return nil, errors.New("user_id not found in context")
	}
	res, err := l.svcCtx.PaymentRpc.CancelPayment(l.ctx, &payment.CancelPaymentReq{
		UserId:  userId,
		OrderSn: req.OrderSn,
	})
	if err != nil {
		return nil, err
	}

	return &types.OrdersCancelResp{
		Message: res.Message,
	}, nil
}

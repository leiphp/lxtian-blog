package payment

import (
	"context"
	"errors"
	"lxtian-blog/rpc/payment/pb/payment"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrdersCloseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 关闭订单
func NewOrdersCloseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrdersCloseLogic {
	return &OrdersCloseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrdersCloseLogic) OrdersClose(req *types.OrdersCloseReq) (resp *types.OrdersCloseResp, err error) {
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		logx.Errorf("OrdersClose user_id not found in context")
		return nil, errors.New("user_id not found in context")
	}
	res, err := l.svcCtx.PaymentRpc.ClosePayment(l.ctx, &payment.ClosePaymentReq{
		UserId:     int64(userId),
		OutTradeNo: req.OutTradeNo,
	})
	if err != nil {
		return nil, err
	}

	return &types.OrdersCloseResp{
		Message: res.Message,
	}, nil
}

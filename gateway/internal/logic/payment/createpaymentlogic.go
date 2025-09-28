package payment

import (
	"context"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"lxtian-blog/rpc/payment/pb/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建支付订单
func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePaymentLogic) CreatePayment(req *types.CreatePaymentReq) (resp *types.CreatePaymentResp, err error) {
	res, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		GoodsId: req.GoodsId,
		UserId:  req.UserId,
		Amount:  req.Amount,
		Subject: req.Subject,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreatePaymentResp{
		PaymentId: res.PaymentId,
	}, nil

	//return
}

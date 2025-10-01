package payment

import (
	"context"
	"errors"
	"lxtian-blog/common/pkg/utils"
	"net/http"

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

func (l *CreatePaymentLogic) CreatePayment(req *types.CreatePaymentReq, r *http.Request) (resp *types.CreatePaymentResp, err error) {
	//从中间件获取用户信息
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		return nil, errors.New("user_id not found in context")
	}
	res, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		GoodsId:   int64(req.GoodsId),
		Quantity:  req.Quantity,
		UserId:    uint64(userId),
		Amount:    req.Amount,
		Subject:   req.Subject,
		ClientIp:  utils.GetClientIp(r),
		PayType:   int64(req.PayType),
		BuyType:   int64(req.BuyType),
		Body:      req.Body,
		NotifyUrl: req.NotifyUrl,
		ReturnUrl: req.ReturnUrl,
		Remark:    req.Remark,
		Timeout:   req.Timeout,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreatePaymentResp{
		PaymentId:  res.PaymentId,
		PayUrl:     res.PayUrl,
		OutTradeNo: res.OutTradeNo,
		OrderSn:    res.OrderSn,
	}, nil

	//return
}

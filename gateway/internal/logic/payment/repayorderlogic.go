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

type RepayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重新支付订单
func NewRepayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RepayOrderLogic {
	return &RepayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RepayOrderLogic) RepayOrder(req *types.RepayOrderReq, r *http.Request) (resp *types.RepayOrderResp, err error) {
	// 从中间件获取用户信息
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		return nil, errors.New("user_id not found in context")
	}

	// 参数验证：order_id 和 out_trade_no 至少提供一个
	if req.OrderId == "" && req.OutTradeNo == "" {
		return nil, errors.New("订单ID或商户订单号至少提供一个")
	}

	// 调用 RPC 服务重新支付订单
	res, err := l.svcCtx.PaymentRpc.RepayOrder(l.ctx, &payment.RepayOrderReq{
		OrderSn:    req.OrderId,
		OutTradeNo: req.OutTradeNo,
		UserId:     uint64(userId),
		ReturnUrl:  req.ReturnUrl,
		NotifyUrl:  req.NotifyUrl,
		ClientIp:   utils.GetClientIp(r),
	})
	if err != nil {
		return nil, err
	}

	return &types.RepayOrderResp{
		PaymentId:  res.PaymentId,
		OutTradeNo: res.OutTradeNo,
		OrderSn:    res.OrderSn,
		PayUrl:     res.PayUrl,
	}, nil
}

package payment

import (
	"context"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/payment/pb/payment"
	"net/http"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DonateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在线捐赠
func NewDonateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DonateLogic {
	return &DonateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DonateLogic) Donate(req *types.DonateReq, r *http.Request) (resp *types.DonateResp, err error) {
	l.Infof("Donate logic: req=%v", req)
	//从中间件获取用户信息
	userId, _ := l.ctx.Value("user_id").(uint)

	res, err := l.svcCtx.PaymentRpc.Donate(l.ctx, &payment.DonateReq{
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

	return &types.DonateResp{
		PaymentId:  res.PaymentId,
		PayUrl:     res.PayUrl,
		OutTradeNo: res.OutTradeNo,
		OrderSn:    res.OrderSn,
	}, nil
}

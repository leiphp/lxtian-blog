package payment

import (
	"context"
	"errors"
	"strings"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"lxtian-blog/rpc/payment/pb/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentNotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付结果异步通知
func NewPaymentNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentNotifyLogic {
	return &PaymentNotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentNotifyLogic) PaymentNotify(req *types.PaymentNotifyReq) (resp *types.PaymentNotifyResp, err error) {
	l.Infof("payment notify req: %+v", req)
	if req == nil {
		return nil, errors.New("request can not be nil")
	}

	notifyData := strings.TrimSpace(req.NotifyData)
	if notifyData == "" {
		return nil, errors.New("notify_data can not be empty")
	}

	// 签名字段可以为空（部分场景可能不要求），这里不做强校验，仅清理空白字符
	sign := strings.TrimSpace(req.Sign)
	signType := strings.TrimSpace(req.SignType)

	clientIP := ""
	if v, ok := l.ctx.Value("client_ip").(string); ok && v != "" {
		clientIP = v
	}

	rpcResp, err := l.svcCtx.PaymentRpc.PaymentNotify(l.ctx, &payment.PaymentNotifyReq{
		NotifyData: notifyData,
		Sign:       sign,
		SignType:   signType,
		ClientIp:   clientIP,
	})
	if err != nil {
		return nil, err
	}

	// 如果回调处理成功，返回 "success" 字符串
	if rpcResp.Success {
		return &types.PaymentNotifyResp{
			Result: "success",
		}, nil
	}

	// 如果处理失败，返回错误信息
	return &types.PaymentNotifyResp{
		Result: rpcResp.Message,
	}, nil
}

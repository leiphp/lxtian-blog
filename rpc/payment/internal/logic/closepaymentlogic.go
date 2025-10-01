package logic

import (
	"context"
	"fmt"
	"lxtian-blog/common/constant"
	paymentSvc "lxtian-blog/common/repository/payment"

	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type ClosePaymentLogic struct {
	*BaseLogic
	paymentService paymentSvc.PaymentOrderRepository
}

func NewClosePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClosePaymentLogic {
	return &ClosePaymentLogic{
		BaseLogic:      NewBaseLogic(ctx, svcCtx),
		paymentService: paymentSvc.NewPaymentOrderRepository(svcCtx.DB),
	}
}

func (l *ClosePaymentLogic) ClosePayment(in *payment.ClosePaymentReq) (*payment.ClosePaymentResp, error) {
	// 参数验证
	if in.PaymentId == "" && in.OrderId == "" && in.OutTradeNo == "" {
		return &payment.ClosePaymentResp{
			Success: false,
			Message: "支付ID、订单ID或商户订单号至少提供一个",
		}, fmt.Errorf("at least one of payment_id, order_id, out_trade_no is required")
	}

	var paymentOrder *model.LxtPaymentOrders
	var err error

	// 根据提供的参数查找支付订单
	if in.PaymentId != "" {
		paymentOrder, err = l.paymentService.GetByPaymentId(l.ctx, in.PaymentId)
	} else if in.OrderId != "" {
		paymentOrder, err = l.paymentService.GetByOrderId(l.ctx, in.OrderId)
	} else if in.OutTradeNo != "" {
		paymentOrder, err = l.paymentService.GetByOutTradeNo(l.ctx, in.OutTradeNo)
	}

	if err != nil {
		l.Errorf("Failed to find payment order: %v", err)
		return &payment.ClosePaymentResp{
			Success: false,
			Message: "支付订单不存在",
		}, fmt.Errorf("payment order not found: %w", err)
	}

	// 检查订单状态是否允许关闭
	if paymentOrder.Status != constant.PaymentStatusPending {
		return &payment.ClosePaymentResp{
			Success: false,
			Message: "订单状态不允许关闭",
		}, fmt.Errorf("order status does not allow close")
	}

	// 调用支付宝API关闭订单
	alipayReq := &alipay.TradeCloseRequest{
		OutTradeNo: paymentOrder.OutTradeNo,
		TradeNo:    paymentOrder.TradeNo,
	}

	alipayResp, err := l.svcCtx.AlipayClient.ClosePayment(alipayReq)
	if err != nil {
		l.Errorf("Failed to close alipay payment: %v", err)
		return &payment.ClosePaymentResp{
			Success: false,
			Message: "关闭支付订单失败",
		}, fmt.Errorf("failed to close alipay payment: %w", err)
	}

	// 更新本地订单状态
	err = l.paymentService.UpdateStatus(l.ctx, paymentOrder.PaymentId, constant.PaymentStatusClosed)
	if err != nil {
		l.Errorf("Failed to update payment status: %v", err)
		// 即使本地更新失败，支付宝那边已经关闭了，所以仍然返回成功
	}

	// 记录日志
	l.Infof("Closed payment order: paymentId=%s, orderSn=%s, outTradeNo=%s",
		paymentOrder.PaymentId, paymentOrder.OrderSn, alipayResp.OutTradeNo)

	return &payment.ClosePaymentResp{
		Success: true,
		Message: "关闭成功",
	}, nil
}

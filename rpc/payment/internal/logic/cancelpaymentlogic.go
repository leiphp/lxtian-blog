package logic

import (
	"context"
	"fmt"

	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/common/pkg/model"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type CancelPaymentLogic struct {
	*BaseLogic
}

func NewCancelPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPaymentLogic {
	return &CancelPaymentLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
	}
}

func (l *CancelPaymentLogic) CancelPayment(in *payment.CancelPaymentReq) (*payment.CancelPaymentResp, error) {
	// 参数验证
	if in.PaymentId == "" && in.OrderId == "" && in.OutTradeNo == "" {
		return &payment.CancelPaymentResp{
			Success: false,
			Message: "支付ID、订单ID或商户订单号至少提供一个",
		}, fmt.Errorf("at least one of payment_id, order_id, out_trade_no is required")
	}

	var paymentOrder *model.PaymentOrder
	var err error

	// 根据提供的参数查找支付订单
	if in.PaymentId != "" {
		paymentOrder, err = l.svcCtx.PaymentModel.FindPaymentOrderByPaymentId(l.ctx, in.PaymentId)
	} else if in.OrderId != "" {
		paymentOrder, err = l.svcCtx.PaymentModel.FindPaymentOrderByOrderId(l.ctx, in.OrderId)
	} else if in.OutTradeNo != "" {
		paymentOrder, err = l.svcCtx.PaymentModel.FindPaymentOrderByOutTradeNo(l.ctx, in.OutTradeNo)
	}

	if err != nil {
		l.Errorf("Failed to find payment order: %v", err)
		return &payment.CancelPaymentResp{
			Success: false,
			Message: "支付订单不存在",
		}, fmt.Errorf("payment order not found: %w", err)
	}

	// 检查订单状态是否允许取消
	if paymentOrder.Status != model.PaymentStatusPending {
		return &payment.CancelPaymentResp{
			Success: false,
			Message: "订单状态不允许取消",
		}, fmt.Errorf("order status does not allow cancel")
	}

	// 调用支付宝API取消订单
	alipayReq := &alipay.TradeCancelRequest{
		OutTradeNo: paymentOrder.OutTradeNo,
		TradeNo:    paymentOrder.TradeNo,
	}

	alipayResp, err := l.svcCtx.AlipayClient.CancelPayment(alipayReq)
	if err != nil {
		l.Errorf("Failed to cancel alipay payment: %v", err)
		return &payment.CancelPaymentResp{
			Success: false,
			Message: "取消支付订单失败",
		}, fmt.Errorf("failed to cancel alipay payment: %w", err)
	}

	// 更新本地订单状态
	err = l.svcCtx.PaymentModel.UpdatePaymentOrderStatus(l.ctx, paymentOrder.PaymentId, model.PaymentStatusCancelled)
	if err != nil {
		l.Errorf("Failed to update payment status: %v", err)
		// 即使本地更新失败，支付宝那边已经取消了，所以仍然返回成功
	}

	// 记录日志
	l.Infof("Cancelled payment order: paymentId=%s, orderId=%s, outTradeNo=%s",
		paymentOrder.PaymentId, paymentOrder.OrderId, alipayResp.OutTradeNo)

	return &payment.CancelPaymentResp{
		Success: true,
		Message: "取消成功",
	}, nil
}

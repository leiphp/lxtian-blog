package logic

import (
	"context"
	"fmt"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/repository/payment_repo"

	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
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
	if in.PaymentId == "" && in.OrderSn == "" && in.OutTradeNo == "" {
		return &payment.CancelPaymentResp{
			Success: false,
			Message: "支付ID、订单ID或商户订单号至少提供一个",
		}, fmt.Errorf("at least one of payment_id, order_id, out_trade_no is required")
	}

	var paymentOrder *model.LxtPaymentOrder
	var err error
	paymentService := payment_repo.NewPaymentOrderRepository(l.svcCtx.DB)
	// 根据提供的参数查找支付订单
	if in.PaymentId != "" {
		paymentOrder, err = paymentService.GetByPaymentId(l.ctx, in.PaymentId)
	} else if in.OrderSn != "" {
		paymentOrder, err = paymentService.GetByOrderSn(l.ctx, in.OrderSn)
	} else if in.OutTradeNo != "" {
		paymentOrder, err = paymentService.GetByOutTradeNo(l.ctx, in.OutTradeNo)
	}

	if err != nil {
		l.Errorf("Failed to find payment order: %v", err)
		return &payment.CancelPaymentResp{
			Success: false,
			Message: "支付订单不存在",
		}, fmt.Errorf("payment order not found: %w", err)
	}

	// 检查订单状态是否允许取消
	if paymentOrder.Status != constant.PaymentStatusPending {
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
	err = paymentService.UpdateStatus(l.ctx, paymentOrder.PaymentID, constant.PaymentStatusCancelled)
	if err != nil {
		l.Errorf("Failed to update payment status: %v", err)
		// 即使本地更新失败，支付宝那边已经取消了，所以仍然返回成功
	}

	// 记录日志
	l.Infof("Cancelled payment order: paymentId=%s, orderSn=%s, outTradeNo=%s",
		paymentOrder.PaymentID, paymentOrder.OrderSn, alipayResp.OutTradeNo)

	return &payment.CancelPaymentResp{
		Success: true,
		Message: "取消成功",
	}, nil
}

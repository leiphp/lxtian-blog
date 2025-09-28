package logic

import (
	"context"
	"fmt"
	"time"

	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/common/pkg/model"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type QueryPaymentLogic struct {
	*BaseLogic
}

func NewQueryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPaymentLogic {
	return &QueryPaymentLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
	}
}

func (l *QueryPaymentLogic) QueryPayment(in *payment.QueryPaymentReq) (*payment.QueryPaymentResp, error) {
	// 参数验证
	if in.PaymentId == "" && in.OrderId == "" && in.OutTradeNo == "" {
		return &payment.QueryPaymentResp{
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
		return &payment.QueryPaymentResp{
			Message: "支付订单不存在",
		}, fmt.Errorf("payment order not found: %w", err)
	}

	// 如果订单状态已经是最终状态，直接返回
	if paymentOrder.Status == model.PaymentStatusPaid ||
		paymentOrder.Status == model.PaymentStatusClosed ||
		paymentOrder.Status == model.PaymentStatusCancelled {
		return l.buildQueryResponse(paymentOrder), nil
	}

	// 调用支付宝API查询最新状态
	alipayReq := &alipay.TradeQueryRequest{
		OutTradeNo: paymentOrder.OutTradeNo,
		TradeNo:    paymentOrder.TradeNo,
	}

	alipayResp, err := l.svcCtx.AlipayClient.QueryPayment(alipayReq)
	if err != nil {
		l.Errorf("Failed to query alipay payment: %v", err)
		// 即使支付宝查询失败，也返回本地数据库的信息
		return l.buildQueryResponse(paymentOrder), nil
	}

	// 更新本地订单状态
	err = l.updatePaymentStatus(paymentOrder, alipayResp)
	if err != nil {
		l.Errorf("Failed to update payment status: %v", err)
	}

	// 重新查询更新后的订单信息
	updatedOrder, err := l.svcCtx.PaymentModel.FindPaymentOrderByPaymentId(l.ctx, paymentOrder.PaymentId)
	if err != nil {
		// 如果查询失败，使用原始数据
		updatedOrder = paymentOrder
	}

	return l.buildQueryResponse(updatedOrder), nil
}

// 更新支付状态
func (l *QueryPaymentLogic) updatePaymentStatus(paymentOrder *model.PaymentOrder, alipayResp *alipay.TradeQueryResponse) error {
	// 根据支付宝返回的交易状态更新本地订单状态
	switch alipayResp.TradeStatus {
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		// 支付成功
		if paymentOrder.Status != model.PaymentStatusPaid {
			// 解析支付时间
			var gmtPayment *time.Time
			if alipayResp.GmtPayment != "" {
				if t, err := time.Parse("2006-01-02 15:04:05", alipayResp.GmtPayment); err == nil {
					gmtPayment = &t
				}
			}

			err := l.svcCtx.PaymentModel.UpdatePaymentOrderTradeInfo(
				l.ctx,
				paymentOrder.PaymentId,
				alipayResp.TradeNo,
				alipayResp.TradeStatus,
				alipayResp.BuyerUserId,
				alipayResp.BuyerLogonId,
				alipayResp.ReceiptAmount,
				gmtPayment,
			)
			if err != nil {
				return err
			}
		}
	case "TRADE_CLOSED":
		// 交易关闭
		if paymentOrder.Status != model.PaymentStatusClosed {
			err := l.svcCtx.PaymentModel.UpdatePaymentOrderStatus(l.ctx, paymentOrder.PaymentId, model.PaymentStatusClosed)
			if err != nil {
				return err
			}
		}
	case "WAIT_BUYER_PAY":
		// 等待买家付款，保持待支付状态
		// 不需要更新状态
	}

	return nil
}

// 构建查询响应
func (l *QueryPaymentLogic) buildQueryResponse(paymentOrder *model.PaymentOrder) *payment.QueryPaymentResp {
	resp := &payment.QueryPaymentResp{
		PaymentId:     paymentOrder.PaymentId,
		OrderId:       paymentOrder.OrderId,
		OutTradeNo:    paymentOrder.OutTradeNo,
		TradeNo:       paymentOrder.TradeNo,
		TradeStatus:   paymentOrder.TradeStatus,
		TotalAmount:   paymentOrder.Amount,
		ReceiptAmount: paymentOrder.ReceiptAmount,
		BuyerUserId:   paymentOrder.BuyerUserId,
		BuyerLogonId:  paymentOrder.BuyerLogonId,
		Message:       "查询成功",
	}

	// 设置支付时间
	if paymentOrder.GmtPayment != nil {
		resp.GmtPayment = paymentOrder.GmtPayment.Format("2006-01-02 15:04:05")
	}

	// 设置关闭时间
	if paymentOrder.GmtClose != nil {
		resp.GmtClose = paymentOrder.GmtClose.Format("2006-01-02 15:04:05")
	}

	return resp
}

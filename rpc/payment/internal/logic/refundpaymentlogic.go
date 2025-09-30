package logic

import (
	"context"
	"database/sql"
	"fmt"
	"lxtian-blog/common/constant"
	paymentSvc "lxtian-blog/common/repository/payment"
	"time"

	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type RefundPaymentLogic struct {
	*BaseLogic
	repo paymentSvc.LxtPaymentRefundsRepo
}

func NewRefundPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundPaymentLogic {
	return &RefundPaymentLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
		repo:      paymentSvc.NewLxtPaymentRefundsRepo(svcCtx.DB),
	}
}

func (l *RefundPaymentLogic) RefundPayment(in *payment.RefundPaymentReq) (*payment.RefundPaymentResp, error) {
	// 参数验证
	if in.PaymentId == "" && in.OrderId == "" && in.OutTradeNo == "" {
		return &payment.RefundPaymentResp{
			Message: "支付ID、订单ID或商户订单号至少提供一个",
		}, fmt.Errorf("at least one of payment_id, order_id, out_trade_no is required")
	}

	if in.RefundAmount <= 0 {
		return &payment.RefundPaymentResp{
			Message: "退款金额必须大于0",
		}, fmt.Errorf("refund_amount must be greater than 0")
	}

	var paymentOrder *model.LxtPaymentOrders
	var err error

	// 根据提供的参数查找支付订单
	if in.PaymentId != "" {
		paymentOrder, err = l.repo.FindPaymentOrderByPaymentId(l.ctx, in.PaymentId)
	} else if in.OrderId != "" {
		paymentOrder, err = l.repo.FindPaymentOrderByOrderId(l.ctx, in.OrderId)
	} else if in.OutTradeNo != "" {
		paymentOrder, err = l.repo.FindPaymentOrderByOutTradeNo(l.ctx, in.OutTradeNo)
	}

	if err != nil {
		l.Errorf("Failed to find payment order: %v", err)
		return &payment.RefundPaymentResp{
			Message: "支付订单不存在",
		}, fmt.Errorf("payment order not found: %w", err)
	}

	// 检查订单状态是否允许退款
	if paymentOrder.Status != constant.PaymentStatusPaid {
		return &payment.RefundPaymentResp{
			Message: "订单状态不允许退款",
		}, fmt.Errorf("order status does not allow refund")
	}

	// 检查退款金额
	if in.RefundAmount > paymentOrder.Amount {
		return &payment.RefundPaymentResp{
			Message: "退款金额不能超过支付金额",
		}, fmt.Errorf("refund amount exceeds payment amount")
	}

	// 生成退款ID和退款单号
	refundId := l.generateRefundId()
	outRequestNo := in.OutRequestNo
	if outRequestNo == "" {
		outRequestNo = fmt.Sprintf("REFUND_%s_%d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
	}

	// 检查是否已有相同的退款单号
	existingRefund, err := l.repo.FindPaymentRefundByOutRequestNo(l.ctx, outRequestNo)
	if err == nil && existingRefund != nil {
		return &payment.RefundPaymentResp{
			Message: "退款单号已存在",
		}, fmt.Errorf("refund request number already exists")
	}

	// 创建退款记录
	paymentRefund := &model.LxtPaymentRefunds{
		RefundId:     refundId,
		PaymentId:    paymentOrder.PaymentId,
		OrderId:      paymentOrder.OrderId,
		OutTradeNo:   paymentOrder.OutTradeNo,
		OutRequestNo: outRequestNo,
		UserId:       paymentOrder.UserId,
		RefundAmount: in.RefundAmount,
		RefundReason: sql.NullString{String: in.RefundReason, Valid: in.RefundReason != ""},
		Status:       constant.RefundStatusPending,
		RefundStatus: sql.NullString{String: "", Valid: false},
	}

	// 保存退款记录到数据库
	_, err = l.repo.InsertPaymentRefund(l.ctx, paymentRefund)
	if err != nil {
		l.Errorf("Failed to insert payment refund: %v", err)
		return &payment.RefundPaymentResp{
			Message: "创建退款记录失败",
		}, fmt.Errorf("failed to insert payment refund: %w", err)
	}

	// 调用支付宝API申请退款
	alipayReq := &alipay.TradeRefundRequest{
		OutTradeNo:   paymentOrder.OutTradeNo,
		TradeNo:      paymentOrder.TradeNo,
		RefundAmount: in.RefundAmount,
		RefundReason: in.RefundReason,
		OutRequestNo: outRequestNo,
	}

	alipayResp, err := l.svcCtx.AlipayClient.RefundPayment(alipayReq)
	if err != nil {
		l.Errorf("Failed to refund alipay payment: %v", err)
		// 更新退款状态为失败
		l.repo.UpdatePaymentRefundStatus(l.ctx, refundId, constant.RefundStatusFailed)
		return &payment.RefundPaymentResp{
			Message: "申请退款失败",
		}, fmt.Errorf("failed to refund alipay payment: %w", err)
	}

	// 更新退款记录
	paymentRefund.RefundStatus = sql.NullString{String: alipayResp.RefundStatus, Valid: alipayResp.RefundStatus != ""}
	paymentRefund.RefundFee = sql.NullFloat64{Float64: alipayResp.RefundFee, Valid: alipayResp.RefundFee > 0}

	// 解析退款时间
	if alipayResp.GmtRefund != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", alipayResp.GmtRefund); err == nil {
			paymentRefund.GmtRefund = sql.NullTime{Time: t, Valid: true}
		}
	}

	// 根据支付宝返回的状态更新本地状态
	switch alipayResp.RefundStatus {
	case "REFUND_SUCCESS":
		paymentRefund.Status = constant.RefundStatusSuccess
	case "REFUND_CLOSED":
		paymentRefund.Status = constant.RefundStatusClosed
	default:
		paymentRefund.Status = constant.RefundStatusFailed
	}

	// 更新退款记录
	err = l.repo.UpdatePaymentRefund(l.ctx, paymentRefund)
	if err != nil {
		l.Errorf("Failed to update payment refund: %v", err)
	}

	// 如果退款成功，更新支付订单状态
	if paymentRefund.Status == constant.RefundStatusSuccess {
		if in.RefundAmount >= paymentOrder.Amount {
			// 全额退款
			l.repo.UpdatePaymentOrderStatus(l.ctx, paymentOrder.PaymentId, constant.PaymentStatusRefunded)
		} else {
			// 部分退款
			l.repo.UpdatePaymentOrderStatus(l.ctx, paymentOrder.PaymentId, constant.PaymentStatusPartialRefunded)
		}
	}

	// 记录日志
	l.Infof("Refund payment: refundId=%s, paymentId=%s, amount=%.2f",
		refundId, paymentOrder.PaymentId, in.RefundAmount)

	resp := &payment.RefundPaymentResp{
		RefundId:     refundId,
		OutRequestNo: outRequestNo,
		RefundAmount: in.RefundAmount,
		Message:      "退款申请成功",
	}

	// 设置退款手续费
	if paymentRefund.RefundFee.Valid {
		resp.RefundFee = paymentRefund.RefundFee.Float64
	}

	// 设置退款状态
	if paymentRefund.RefundStatus.Valid {
		resp.RefundStatus = paymentRefund.RefundStatus.String
	}

	// 设置退款时间
	if paymentRefund.GmtRefund.Valid {
		resp.GmtRefund = paymentRefund.GmtRefund.Time.Format("2006-01-02 15:04:05")
	}

	return resp, nil
}

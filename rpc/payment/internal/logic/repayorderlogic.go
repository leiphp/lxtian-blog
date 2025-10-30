package logic

import (
	"context"
	"fmt"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
	"strconv"

	"gorm.io/gorm"
)

type RepayOrderLogic struct {
	*BaseLogic
}

func NewRepayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RepayOrderLogic {
	return &RepayOrderLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
	}
}

// RepayOrder 重新支付订单（针对已存在的未支付订单）
func (l *RepayOrderLogic) RepayOrder(in *payment.RepayOrderReq) (*payment.RepayOrderResp, error) {
	// 参数验证：order_id 和 out_trade_no 至少提供一个
	if in.OrderId == "" && in.OutTradeNo == "" {
		return nil, fmt.Errorf("订单ID或商户订单号至少提供一个")
	}

	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 1. 查询订单信息
	var paymentOrder model.LxtPaymentOrder
	query := l.svcCtx.DB.WithContext(l.ctx)

	if in.OrderId != "" {
		query = query.Where("order_sn = ?", in.OrderId)
	} else {
		query = query.Where("out_trade_no = ?", in.OutTradeNo)
	}

	err := query.First(&paymentOrder).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("订单不存在")
		}
		l.Errorf("Failed to query payment_repo order: %v", err)
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	// 2. 验证订单是否属于当前用户
	if paymentOrder.UserID != int64(in.UserId) {
		return nil, fmt.Errorf("无权操作此订单")
	}

	// 3. 验证订单状态是否为未支付
	if paymentOrder.Status != constant.PaymentStatusPending {
		return nil, fmt.Errorf("订单状态不是待支付，当前状态：%s", paymentOrder.Status)
	}

	// 4. 生成新的支付ID（每次发起支付都生成新的payment_id，但使用相同的out_trade_no）
	newPaymentId := l.generatePaymentId()

	// 5. 调用支付宝API创建支付订单
	amountStr := strconv.FormatFloat(paymentOrder.Amount, 'f', 2, 64)

	// 超时时间
	timeout := paymentOrder.Timeout
	if timeout == "" {
		timeout = "30m"
	}

	// 使用回调地址，优先使用请求中的，否则使用订单中的
	returnUrl := in.ReturnUrl
	if returnUrl == "" {
		returnUrl = paymentOrder.ReturnURL
	}

	notifyUrl := in.NotifyUrl
	if notifyUrl == "" {
		notifyUrl = paymentOrder.NotifyURL
	}

	alipayReq := &alipay.TradeCreateRequest{
		OutTradeNo:  paymentOrder.OutTradeNo,
		TotalAmount: amountStr,
		Subject:     paymentOrder.Subject,
		Body:        *paymentOrder.Body,
		ProductCode: *paymentOrder.ProductCode,
		Timeout:     timeout,
		ReturnUrl:   returnUrl,
	}

	// 调用支付宝API创建支付URL
	payUrl, err := l.svcCtx.AlipayClient.CreatePayment(alipayReq)
	if err != nil {
		l.Errorf("Failed to create alipay payment_repo: %v", err)
		return nil, fmt.Errorf("创建支付订单失败: %w", err)
	}

	// 6. 更新订单的 payment_id（可选，如果需要记录最新的支付ID）
	err = l.svcCtx.DB.WithContext(l.ctx).Model(&model.LxtPaymentOrder{}).
		Where("id = ?", paymentOrder.ID).
		Update("payment_id", newPaymentId).Error
	if err != nil {
		l.Errorf("Failed to update payment_id: %v", err)
		// 这里不返回错误，因为支付链接已经生成
	}

	// 记录日志
	l.Infof("Repay order: paymentId=%s, orderSn=%s, outTradeNo=%s, amount=%.2f, payUrl=%s",
		newPaymentId, paymentOrder.OrderSn, paymentOrder.OutTradeNo, paymentOrder.Amount, payUrl)

	return &payment.RepayOrderResp{
		PaymentId:  newPaymentId,
		OutTradeNo: paymentOrder.OutTradeNo,
		OrderSn:    paymentOrder.OrderSn,
		PayUrl:     payUrl,
	}, nil
}

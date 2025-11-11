package logic

import (
	"context"
	"fmt"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/common/pkg/utils"
	"strconv"

	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type DonateLogic struct {
	*BaseLogic
}

func NewDonateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DonateLogic {
	return &DonateLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
	}
}

// 创建捐赠订单
func (l *DonateLogic) Donate(in *payment.DonateReq) (*payment.DonateResp, error) {
	// 参数验证
	if in.Amount <= 0 {
		return nil, fmt.Errorf("支付金额必须大于0")
	}

	if in.Subject == "" {
		return nil, fmt.Errorf("订单标题不能为空")
	}

	// 生成订单ID、支付ID和商户订单号
	paymentId := l.generatePaymentId()
	outTradeNo := fmt.Sprintf("NO%d", utils.Snowflake())
	orderSn := fmt.Sprintf("%d", utils.Snowflake())

	// 产品码：电脑网站支付固定使用 FAST_INSTANT_TRADE_PAY
	if in.ProductCode == "" {
		in.ProductCode = "FAST_INSTANT_TRADE_PAY"
	}

	if in.Timeout == "" {
		in.Timeout = "30m"
	}
	// 设置默认值（需要重新生成protobuf后启用）
	if in.PayType == 0 {
		in.PayType = 1 // 默认直接消费
	}
	// 设置默认值（需要重新生成protobuf后启用）
	if in.BuyType == 0 {
		in.PayType = 3 // 默认直接消费
	}

	// 2. 创建支付订单记录（txy_orders表）
	paymentOrder := &model.TxyOrder{
		PaymentID:  paymentId,
		OrderSn:    orderSn,
		OutTradeNo: outTradeNo,
		UserID:     int32(in.UserId),
		PayMoney:   in.Amount,
		GoodsName:  in.Subject,
		Status:     0,
		PayType:    int32(in.PayType),
		IP:         in.ClientIp,
	}

	// 使用GORM保存支付订单到数据库
	err := l.svcCtx.DB.WithContext(l.ctx).Create(paymentOrder).Error
	if err != nil {
		l.Errorf("Failed to insert payment_repo order: %v", err)
		return nil, fmt.Errorf("创建支付订单失败: %w", err)
	}

	// 3. 调用支付宝API创建支付订单
	// 金额转换为字符串，保留2位小数
	amountStr := strconv.FormatFloat(in.Amount, 'f', 2, 64)

	// 超时时间格式：30m, 1h, 1d 等，默认30m
	timeout := in.Timeout
	if timeout == "" || !isValidTimeout(timeout) {
		timeout = "30m"
	}

	alipayReq := &alipay.TradeCreateRequest{
		OutTradeNo:  outTradeNo,
		TotalAmount: amountStr,
		Subject:     in.Subject,
		Body:        in.Body,
		ProductCode: in.ProductCode,
		Timeout:     timeout,
		ReturnUrl:   in.ReturnUrl,
	}

	// 调用支付宝API创建支付URL
	payUrl, err := l.svcCtx.AlipayClient.CreatePayment(alipayReq)
	if err != nil {
		l.Errorf("Failed to create alipay payment_repo: %v", err)
		// 使用GORM更新订单状态为失败
		l.svcCtx.DB.WithContext(l.ctx).Model(&model.LxtPaymentOrder{}).
			Where("payment_id = ?", paymentId).
			Update("status", constant.VerifyStatusFailed)
		return nil, fmt.Errorf("创建支付订单失败: %w", err)
	}

	// 记录日志
	l.Infof("Created payment_repo order: paymentId=%s, orderSn=%s, outTradeNo=%s, amount=%.2f, payUrl=%s",
		paymentId, orderSn, outTradeNo, in.Amount, payUrl)

	return &payment.DonateResp{
		PaymentId:  paymentId,
		OutTradeNo: outTradeNo,
		OrderSn:    orderSn,
		PayUrl:     payUrl,
	}, nil
}

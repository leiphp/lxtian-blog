package logic

import (
	"context"
	"database/sql"
	"fmt"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
	"strconv"
)

type CreatePaymentLogic struct {
	*BaseLogic
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
	}
}

func (l *CreatePaymentLogic) CreatePayment(in *payment.CreatePaymentReq) (*payment.CreatePaymentResp, error) {
	// 参数验证
	if in.Amount <= 0 {
		return nil, fmt.Errorf("支付金额必须大于0")
	}

	if in.Subject == "" {
		return nil, fmt.Errorf("订单标题不能为空")
	}

	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
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
	if in.Quantity == 0 {
		in.Quantity = 1 // 默认数量为1
	}

	// 2. 创建支付订单记录（lxt_payment_orders表）
	paymentOrder := &model.LxtPaymentOrders{
		GoodsId:     in.GoodsId,
		PaymentId:   paymentId,
		OrderSn:     orderSn,
		OutTradeNo:  outTradeNo,
		UserId:      int64(in.UserId),
		Amount:      in.Amount,
		Subject:     in.Subject,
		Status:      constant.PaymentStatusPending,
		ProductCode: in.ProductCode,
		PayType:     in.PayType,
		BuyType:     in.BuyType,
		// 处理可空字段
		Body: sql.NullString{
			String: in.Body,
			Valid:  in.Body != "",
		},
		ReturnUrl: in.ReturnUrl,
		NotifyUrl: in.NotifyUrl,
		Timeout:   in.Timeout,
		ClientIp:  in.ClientIp,
	}

	// 使用GORM保存支付订单到数据库
	err := l.svcCtx.DB.WithContext(l.ctx).Create(paymentOrder).Error
	if err != nil {
		l.Errorf("Failed to insert payment order: %v", err)
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
		l.Errorf("Failed to create alipay payment: %v", err)
		// 使用GORM更新订单状态为失败
		l.svcCtx.DB.WithContext(l.ctx).Model(&model.LxtPaymentOrders{}).
			Where("payment_id = ?", paymentId).
			Update("status", constant.VerifyStatusFailed)
		return nil, fmt.Errorf("创建支付订单失败: %w", err)
	}

	// 记录日志
	l.Infof("Created payment order: paymentId=%s, orderSn=%s, outTradeNo=%s, amount=%.2f, payUrl=%s",
		paymentId, orderSn, outTradeNo, in.Amount, payUrl)

	return &payment.CreatePaymentResp{
		PaymentId:  paymentId,
		OutTradeNo: outTradeNo,
		OrderSn:    orderSn,
		PayUrl:     payUrl,
	}, nil
}

// 验证超时时间格式是否有效
func isValidTimeout(timeout string) bool {
	// 支持的格式：30m, 1h, 1d, 1c（c表示天）
	if len(timeout) < 2 {
		return false
	}

	// 检查是否以数字开头，以 m/h/d/c 结尾
	lastChar := timeout[len(timeout)-1]
	return (lastChar == 'm' || lastChar == 'h' || lastChar == 'd' || lastChar == 'c') &&
		timeout[:len(timeout)-1] != "" &&
		isNumeric(timeout[:len(timeout)-1])
}

// 检查字符串是否为数字
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

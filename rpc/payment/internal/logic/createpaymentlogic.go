package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/rpc/payment/internal/svc"
	"strconv"
	"time"

	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
	"lxtian-blog/rpc/payment/pb/payment"
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
		return &payment.CreatePaymentResp{
			Message: "支付金额必须大于0",
		}, fmt.Errorf("amount must be greater than 0")
	}

	if in.Subject == "" {
		return &payment.CreatePaymentResp{
			Message: "订单标题不能为空",
		}, fmt.Errorf("subject is required")
	}

	if in.UserId == 0 {
		return &payment.CreatePaymentResp{
			Message: "用户ID不能为空",
		}, fmt.Errorf("user_id is required")
	}

	// 生成订单ID、支付ID和商户订单号
	orderId := l.generateOrderId()
	paymentId := l.generatePaymentId()
	outTradeNo := fmt.Sprintf("ORDER_%s_%d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000)
	orderSn := fmt.Sprintf("SN%s", time.Now().Format("20060102150405"))

	// 设置默认值
	if in.ProductCode == "" {
		in.ProductCode = "FAST_INSTANT_TRADE_PAY"
	}
	if in.Timeout == "" {
		in.Timeout = "30m"
	}
	// 设置默认值（需要重新生成protobuf后启用）
	// if in.PayType == 0 {
	//     in.PayType = 3 // 默认直接消费
	// }
	// if in.GoodsName == "" {
	//     in.GoodsName = in.Subject
	// }
	// if in.Quantity == 0 {
	//     in.Quantity = 1 // 默认数量为1
	// }

	// 临时默认值
	payType := int64(3) // 默认直接消费
	goodsName := in.Subject
	quantity := uint32(1)
	goodsId := uint64(0)
	remark := ""

	// 1. 先创建业务订单（txy_order表）
	businessOrder := &mysql.TxyOrder{
		OutTradeNo: outTradeNo,
		OrderSn:    orderSn,
		PayMoney:   in.Amount,
		PayType:    payType,
		UserId:     in.UserId,
		Status:     0, // 0待支付
		GoodsName:  goodsName,
		Ctime:      time.Now().Unix(),
		Remark:     fmt.Sprintf("商品ID:%d, 数量:%d, %s", goodsId, quantity, remark),
	}

	// 临时注释，避免未使用变量错误
	_ = businessOrder

	// 这里需要添加业务订单的插入逻辑
	// 由于当前没有业务订单的模型，我们先用注释标记
	// _, err := l.svcCtx.BusinessOrderModel.Insert(l.ctx, businessOrder)
	// if err != nil {
	//     l.Errorf("Failed to insert business order: %v", err)
	//     return &payment.CreatePaymentResp{
	//         Message: "创建业务订单失败",
	//     }, fmt.Errorf("failed to insert business order: %w", err)
	// }

	// 2. 创建支付订单记录（payment_orders表）
	paymentOrder := &model.LxtPaymentOrders{
		PaymentId:     paymentId,
		OrderId:       orderId,
		OutTradeNo:    outTradeNo,
		UserId:        int64(in.UserId),
		Amount:        in.Amount,
		Subject:       in.Subject,
		Body:          in.Body,
		Status:        constant.PaymentStatusPending,
		ProductCode:   in.ProductCode,
		ReturnUrl:     in.ReturnUrl,
		NotifyUrl:     in.NotifyUrl,
		Timeout:       in.Timeout,
		ReceiptAmount: "0.00", // 初始化为0.00，避免decimal字段错误
	}

	// 使用GORM保存支付订单到数据库
	err := l.svcCtx.DB.WithContext(l.ctx).Create(paymentOrder).Error
	if err != nil {
		l.Errorf("Failed to insert payment order: %v", err)
		return &payment.CreatePaymentResp{
			Message: "创建支付订单失败",
		}, fmt.Errorf("failed to insert payment order: %w", err)
	}

	// 3. 调用支付宝API创建支付订单
	// 金额转换为字符串，保留2位小数
	amountStr := strconv.FormatFloat(in.Amount, 'f', 2, 64)

	// 产品码：电脑网站支付固定使用 FAST_INSTANT_TRADE_PAY
	productCode := in.ProductCode
	if productCode == "" || productCode == "Lorem" {
		productCode = "FAST_INSTANT_TRADE_PAY"
	}

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
		ProductCode: productCode,
		Timeout:     timeout,
		ReturnUrl:   in.ReturnUrl,
	}

	alipayResp, err := l.svcCtx.AlipayClient.CreatePayment(alipayReq)
	fmt.Println("alipayResp:", alipayResp)
	if err != nil {
		l.Errorf("Failed to create alipay payment: %v", err)
		// 使用GORM更新订单状态为失败
		l.svcCtx.DB.WithContext(l.ctx).Model(&model.LxtPaymentOrders{}).
			Where("payment_id = ?", paymentId).
			Update("status", constant.VerifyStatusFailed)
		return &payment.CreatePaymentResp{
			Message: "创建支付订单失败",
		}, fmt.Errorf("failed to create alipay payment: %w", err)
	}

	// 构建支付链接
	payUrl := l.buildPaymentUrl(alipayReq)

	// 构建表单数据
	formData := l.buildFormData(alipayReq)

	// 记录日志
	l.Infof("Created payment order: paymentId=%s, orderId=%s, amount=%.2f",
		paymentId, orderId, in.Amount)

	return &payment.CreatePaymentResp{
		PaymentId: paymentId,
		PayUrl:    payUrl,
		QrCode:    alipayResp.QrCode,
		FormData:  formData,
		Message:   "订单创建成功",
	}, nil
}

// 构建支付链接
func (l *CreatePaymentLogic) buildPaymentUrl(req *alipay.TradeCreateRequest) string {
	// 这里可以构建跳转到支付宝的URL
	// 实际实现中需要根据支付宝的文档来构建
	return fmt.Sprintf("https://openapi.alipay.com/gateway.do?out_trade_no=%s&total_amount=%s&subject=%s",
		req.OutTradeNo, req.TotalAmount, req.Subject)
}

// 构建表单数据
func (l *CreatePaymentLogic) buildFormData(req *alipay.TradeCreateRequest) string {
	// 构建支付宝支付表单数据
	formData := map[string]interface{}{
		"out_trade_no":    req.OutTradeNo,
		"total_amount":    req.TotalAmount,
		"subject":         req.Subject,
		"body":            req.Body,
		"product_code":    req.ProductCode,
		"timeout_express": req.Timeout,
	}

	if req.ReturnUrl != "" {
		formData["return_url"] = req.ReturnUrl
	}

	data, _ := json.Marshal(formData)
	return string(data)
}

// 生成订单ID
func (l *CreatePaymentLogic) generateOrderId() string {
	return fmt.Sprintf("ORDER_%d_%d", time.Now().Unix(), time.Now().UnixNano()%100000)
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

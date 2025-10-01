package alipay

import (
	"fmt"
	"net/http"
)

// ExampleCreatePayment 展示如何创建支付订单（电脑网站支付）
// 参考: https://github.com/leiphp/alipay/tree/master
func ExampleCreatePayment(client *AlipayClient, writer http.ResponseWriter, request *http.Request) {
	// 构建支付请求
	var p = &TradeCreateRequest{
		OutTradeNo:  fmt.Sprintf("ORDER_%d", 123456), // 商户订单号
		TotalAmount: "0.01",                          // 支付金额
		Subject:     "测试订单",                          // 订单标题
		Body:        "这是一个测试订单",                      // 订单描述
		ProductCode: "FAST_INSTANT_TRADE_PAY",        // 产品码（电脑网站支付固定值）
		ReturnUrl:   "http://yourdomain.com/return",  // 支付成功后跳转地址
	}

	// 创建支付URL
	payURL, err := client.CreatePayment(p)
	if err != nil {
		// 错误处理
		http.Error(writer, "创建支付失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 重定向到支付宝支付页面
	fmt.Println("支付URL:", payURL)
	http.Redirect(writer, request, payURL, http.StatusTemporaryRedirect)
}

// ExampleCreatePaymentSimple 更简单的示例
func ExampleCreatePaymentSimple() {
	// 1. 创建支付宝客户端配置
	config := &AlipayConfig{
		AppId:           "your_app_id",
		AppPrivateKey:   "your_private_key",
		AlipayPublicKey: "alipay_public_key",
		GatewayUrl:      "https://openapi.alipaydev.com/gateway.do", // 沙箱环境
		NotifyUrl:       "http://yourdomain.com/notify",             // 异步通知地址
		SignType:        "RSA2",
		Charset:         "utf-8",
		Format:          "JSON",
		Version:         "1.0",
	}

	// 2. 创建客户端
	client, err := NewAlipayClient(config)
	if err != nil {
		fmt.Println("创建客户端失败:", err)
		return
	}

	// 3. 创建支付请求
	req := &TradeCreateRequest{
		OutTradeNo:  "ORDER_20240101_001",
		TotalAmount: "99.99",
		Subject:     "购买商品",
		ProductCode: "FAST_INSTANT_TRADE_PAY",
		ReturnUrl:   "http://yourdomain.com/return",
	}

	// 4. 生成支付URL
	payURL, err := client.CreatePayment(req)
	if err != nil {
		fmt.Println("创建支付失败:", err)
		return
	}

	// 5. 使用支付URL
	fmt.Println("支付URL:", payURL)
	// 可以将 payURL 返回给前端，让用户跳转到支付页面
	// 或者在服务端直接重定向: http.Redirect(w, r, payURL, http.StatusTemporaryRedirect)
}

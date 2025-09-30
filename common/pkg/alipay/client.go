package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AlipayClient 支付宝客户端
type AlipayClient struct {
	config *AlipayConfig
	client *http.Client
}

// NewAlipayClient 创建支付宝客户端
func NewAlipayClient(config *AlipayConfig) (*AlipayClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &AlipayClient{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// BaseRequest 基础请求结构
type BaseRequest struct {
	AppId      string `json:"app_id"`
	Method     string `json:"method"`
	Format     string `json:"format"`
	Charset    string `json:"charset"`
	SignType   string `json:"sign_type"`
	Sign       string `json:"sign"`
	Timestamp  string `json:"timestamp"`
	Version    string `json:"version"`
	NotifyUrl  string `json:"notify_url,omitempty"`
	ReturnUrl  string `json:"return_url,omitempty"`
	BizContent string `json:"biz_content"`
}

// BaseResponse 基础响应结构
type BaseResponse struct {
	Code    string          `json:"code"`
	Msg     string          `json:"msg"`
	SubCode string          `json:"sub_code,omitempty"`
	SubMsg  string          `json:"sub_msg,omitempty"`
	Sign    string          `json:"sign,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// TradeCreateRequest 创建支付订单请求
type TradeCreateRequest struct {
	OutTradeNo  string `json:"out_trade_no"`              // 商户订单号
	TotalAmount string `json:"total_amount"`              // 订单总金额（字符串格式，如"88.88"）
	Subject     string `json:"subject"`                   // 订单标题
	Body        string `json:"body,omitempty"`            // 订单描述
	ProductCode string `json:"product_code"`              // 产品码（固定值FAST_INSTANT_TRADE_PAY）
	Timeout     string `json:"timeout_express,omitempty"` // 订单超时时间
	ReturnUrl   string `json:"return_url,omitempty"`      // 支付成功跳转地址
}

// TradeCreateResponse 创建支付订单响应
type TradeCreateResponse struct {
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
	QrCode     string `json:"qr_code"`      // 二维码内容
}

// TradeQueryRequest 查询支付订单请求
type TradeQueryRequest struct {
	OutTradeNo string `json:"out_trade_no"`       // 商户订单号
	TradeNo    string `json:"trade_no,omitempty"` // 支付宝交易号
}

// TradeQueryResponse 查询支付订单响应
type TradeQueryResponse struct {
	OutTradeNo    string  `json:"out_trade_no"`   // 商户订单号
	TradeNo       string  `json:"trade_no"`       // 支付宝交易号
	TradeStatus   string  `json:"trade_status"`   // 交易状态
	TotalAmount   float64 `json:"total_amount"`   // 交易金额
	ReceiptAmount float64 `json:"receipt_amount"` // 实收金额
	BuyerUserId   string  `json:"buyer_user_id"`  // 买家支付宝用户ID
	BuyerLogonId  string  `json:"buyer_logon_id"` // 买家支付宝账号
	GmtPayment    string  `json:"gmt_payment"`    // 支付时间
	GmtClose      string  `json:"gmt_close"`      // 交易关闭时间
}

// TradeRefundRequest 退款请求
type TradeRefundRequest struct {
	OutTradeNo   string  `json:"out_trade_no"`             // 商户订单号
	TradeNo      string  `json:"trade_no,omitempty"`       // 支付宝交易号
	RefundAmount float64 `json:"refund_amount"`            // 退款金额
	RefundReason string  `json:"refund_reason,omitempty"`  // 退款原因
	OutRequestNo string  `json:"out_request_no,omitempty"` // 退款单号
}

// TradeRefundResponse 退款响应
type TradeRefundResponse struct {
	OutTradeNo   string  `json:"out_trade_no"`   // 商户订单号
	OutRequestNo string  `json:"out_request_no"` // 退款单号
	RefundAmount float64 `json:"refund_amount"`  // 退款金额
	RefundFee    float64 `json:"refund_fee"`     // 退款手续费
	RefundStatus string  `json:"refund_status"`  // 退款状态
	GmtRefund    string  `json:"gmt_refund"`     // 退款时间
}

// TradeCloseRequest 关闭订单请求
type TradeCloseRequest struct {
	OutTradeNo string `json:"out_trade_no"`       // 商户订单号
	TradeNo    string `json:"trade_no,omitempty"` // 支付宝交易号
}

// TradeCloseResponse 关闭订单响应
type TradeCloseResponse struct {
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
}

// TradeCancelRequest 取消订单请求
type TradeCancelRequest struct {
	OutTradeNo string `json:"out_trade_no"`       // 商户订单号
	TradeNo    string `json:"trade_no,omitempty"` // 支付宝交易号
}

// TradeCancelResponse 取消订单响应
type TradeCancelResponse struct {
	OutTradeNo string `json:"out_trade_no"` // 商户订单号
}

// CreatePayment 创建支付订单（电脑网站支付）
func (c *AlipayClient) CreatePayment(req *TradeCreateRequest) (*TradeCreateResponse, error) {
	bizContent, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal biz_content: %w", err)
	}

	bizContentStr := string(bizContent)

	// 构建请求参数（除了biz_content，其他参数都放URL中）
	params := url.Values{}
	params.Set("app_id", c.config.AppId)
	params.Set("method", "alipay.trade.page.pay")
	params.Set("format", c.config.Format)
	params.Set("charset", c.config.Charset)
	params.Set("sign_type", c.config.SignType)
	params.Set("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	params.Set("version", c.config.Version)
	params.Set("biz_content", bizContentStr) // 先加入用于签名

	if c.config.NotifyUrl != "" {
		params.Set("notify_url", c.config.NotifyUrl)
	}

	if req.ReturnUrl != "" {
		params.Set("return_url", req.ReturnUrl)
	}

	// 签名（包含biz_content）
	sign, err := c.signParams(params)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}
	params.Set("sign", sign)

	// 移除biz_content，它将放在表单中
	params.Del("biz_content")

	// 构建action URL（不含biz_content）
	actionUrl := c.config.GatewayUrl + "?" + params.Encode()

	// HTML转义biz_content的双引号
	bizContentEscaped := strings.ReplaceAll(bizContentStr, `"`, `&quot;`)

	// 构建HTML表单
	var formBuilder strings.Builder
	formBuilder.WriteString("<form name=\"punchout_form\" method=\"post\" action=\"")
	formBuilder.WriteString(actionUrl)
	formBuilder.WriteString("\">\n")
	formBuilder.WriteString("<input type=\"hidden\" name=\"biz_content\" value=\"")
	formBuilder.WriteString(bizContentEscaped)
	formBuilder.WriteString("\">\n")
	formBuilder.WriteString("<input type=\"submit\" value=\"立即支付\" style=\"display:none\">\n")
	formBuilder.WriteString("</form>\n")
	formBuilder.WriteString("<script>document.forms[0].submit();</script>")

	formHtml := formBuilder.String()

	fmt.Println("=== 生成的HTML表单 ===")
	fmt.Println(formHtml)
	fmt.Println("======================")

	return &TradeCreateResponse{
		OutTradeNo: req.OutTradeNo,
		QrCode:     formHtml, // 返回HTML表单
	}, nil
}

// CreatePaymentQrCode 创建支付订单（当面付-二维码）
func (c *AlipayClient) CreatePaymentQrCode(req *TradeCreateRequest) (*TradeCreateResponse, error) {
	bizContent, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal biz_content: %w", err)
	}

	// 使用当面付接口
	response, err := c.call("alipay.trade.precreate", string(bizContent))
	if err != nil {
		return nil, err
	}

	var result TradeCreateResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// QueryPayment 查询支付订单
func (c *AlipayClient) QueryPayment(req *TradeQueryRequest) (*TradeQueryResponse, error) {
	bizContent, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal biz_content: %w", err)
	}

	response, err := c.call("alipay.trade.query", string(bizContent))
	if err != nil {
		return nil, err
	}

	var result TradeQueryResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// RefundPayment 申请退款
func (c *AlipayClient) RefundPayment(req *TradeRefundRequest) (*TradeRefundResponse, error) {
	bizContent, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal biz_content: %w", err)
	}

	response, err := c.call("alipay.trade.refund", string(bizContent))
	if err != nil {
		return nil, err
	}

	var result TradeRefundResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// ClosePayment 关闭支付订单
func (c *AlipayClient) ClosePayment(req *TradeCloseRequest) (*TradeCloseResponse, error) {
	bizContent, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal biz_content: %w", err)
	}

	response, err := c.call("alipay.trade.close", string(bizContent))
	if err != nil {
		return nil, err
	}

	var result TradeCloseResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// CancelPayment 取消支付订单
func (c *AlipayClient) CancelPayment(req *TradeCancelRequest) (*TradeCancelResponse, error) {
	bizContent, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal biz_content: %w", err)
	}

	response, err := c.call("alipay.trade.cancel", string(bizContent))
	if err != nil {
		return nil, err
	}

	var result TradeCancelResponse
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// call 调用支付宝API
func (c *AlipayClient) call(method, bizContent string) (json.RawMessage, error) {
	// 构建请求参数
	baseReq := &BaseRequest{
		AppId:      c.config.AppId,
		Method:     method,
		Format:     c.config.Format,
		Charset:    c.config.Charset,
		SignType:   c.config.SignType,
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Version:    c.config.Version,
		BizContent: bizContent,
	}

	if c.config.NotifyUrl != "" {
		baseReq.NotifyUrl = c.config.NotifyUrl
	}

	// 签名
	sign, err := c.sign(baseReq)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}
	baseReq.Sign = sign

	// 构建请求URL
	values := url.Values{}
	values.Set("app_id", baseReq.AppId)
	values.Set("method", baseReq.Method)
	values.Set("format", baseReq.Format)
	values.Set("charset", baseReq.Charset)
	values.Set("sign_type", baseReq.SignType)
	values.Set("sign", baseReq.Sign)
	values.Set("timestamp", baseReq.Timestamp)
	values.Set("version", baseReq.Version)
	if baseReq.NotifyUrl != "" {
		values.Set("notify_url", baseReq.NotifyUrl)
	}
	values.Set("biz_content", baseReq.BizContent)

	// 发送请求
	resp, err := c.client.PostForm(c.config.GatewayUrl, values)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 打印原始响应用于调试
	fmt.Printf("Alipay Response Status: %d\n", resp.StatusCode)
	fmt.Printf("Alipay Response Body: %s\n", string(body))

	// 支付宝响应格式：{"alipay_xxx_response": {...}, "sign": "..."}
	// 需要先解析外层结构
	var rawResp map[string]json.RawMessage
	if err := json.Unmarshal(body, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(body))
	}

	// 查找响应节点（通常是 method_response 格式）
	var responseData json.RawMessage
	for key, value := range rawResp {
		if strings.HasSuffix(key, "_response") {
			responseData = value
			break
		}
	}

	if responseData == nil {
		return nil, fmt.Errorf("no response node found in: %s", string(body))
	}

	// 解析响应数据
	var baseResp BaseResponse
	if err := json.Unmarshal(responseData, &baseResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response data: %w, data: %s", err, string(responseData))
	}

	// 检查响应状态
	if baseResp.Code != "10000" {
		// 详细的错误提示
		errMsg := fmt.Sprintf("支付宝错误: code=%s, msg=%s, sub_code=%s, sub_msg=%s",
			baseResp.Code, baseResp.Msg, baseResp.SubCode, baseResp.SubMsg)

		// 针对常见错误提供解决建议
		if baseResp.SubCode == "ACQ.ACCESS_FORBIDDEN" {
			errMsg += "\n【原因】应用没有权限调用此接口"
			errMsg += "\n【解决】1. 登录 https://open.alipay.com 检查应用状态"
			errMsg += "\n       2. 确认已签约'当面付'产品且状态为'已生效'"
			errMsg += "\n       3. 检查应用是否已上线且审核通过"
		}

		return nil, fmt.Errorf(errMsg)
	}

	return responseData, nil
}

// signParams 对 url.Values 参数进行签名
func (c *AlipayClient) signParams(params url.Values) (string, error) {
	// 排序参数key
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "sign" { // 排除sign参数本身
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 构建待签名字符串 key1=value1&key2=value2
	var signStr strings.Builder
	for i, k := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(k)
		signStr.WriteString("=")
		signStr.WriteString(params.Get(k))
	}

	signString := signStr.String()
	fmt.Println("Sign String:", signString)

	// RSA2签名
	privateKeyPEM := c.formatPrivateKey(c.config.AppPrivateKey)

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	hashed := sha256.Sum256([]byte(signString))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// sign 对请求进行签名（兼容旧方法）
func (c *AlipayClient) sign(req *BaseRequest) (string, error) {
	// 构建参数
	params := url.Values{}
	params.Set("app_id", req.AppId)
	params.Set("method", req.Method)
	params.Set("format", req.Format)
	params.Set("charset", req.Charset)
	params.Set("sign_type", req.SignType)
	params.Set("timestamp", req.Timestamp)
	params.Set("version", req.Version)
	params.Set("biz_content", req.BizContent)

	if req.NotifyUrl != "" {
		params.Set("notify_url", req.NotifyUrl)
	}

	if req.ReturnUrl != "" {
		params.Set("return_url", req.ReturnUrl)
	}

	return c.signParams(params)
}

// formatPrivateKey 格式化私钥，自动添加 PEM 头尾
func (c *AlipayClient) formatPrivateKey(key string) string {
	key = strings.TrimSpace(key)

	// 如果已经有 PEM 头尾，直接返回
	if strings.HasPrefix(key, "-----BEGIN") {
		return key
	}

	// 移除所有空格和换行符
	key = strings.ReplaceAll(key, " ", "")
	key = strings.ReplaceAll(key, "\n", "")
	key = strings.ReplaceAll(key, "\r", "")

	// 添加 PEM 头尾，并每64个字符换行
	var formatted strings.Builder
	formatted.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")

	for i := 0; i < len(key); i += 64 {
		end := i + 64
		if end > len(key) {
			end = len(key)
		}
		formatted.WriteString(key[i:end])
		formatted.WriteString("\n")
	}

	formatted.WriteString("-----END RSA PRIVATE KEY-----")
	return formatted.String()
}

// VerifySign 验证签名
func (c *AlipayClient) VerifySign(data, sign string) error {
	// 格式化公钥（自动添加 PEM 头尾，如果没有的话）
	publicKeyPEM := c.formatPublicKey(c.config.AlipayPublicKey)

	// 解析支付宝公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return fmt.Errorf("failed to decode alipay public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse alipay public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("invalid public key type")
	}

	// 解码签名
	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// 验证签名
	hashed := sha256.Sum256([]byte(data))
	return rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signature)
}

// formatPublicKey 格式化公钥，自动添加 PEM 头尾
func (c *AlipayClient) formatPublicKey(key string) string {
	key = strings.TrimSpace(key)

	// 如果已经有 PEM 头尾，直接返回
	if strings.HasPrefix(key, "-----BEGIN") {
		return key
	}

	// 移除所有空格和换行符
	key = strings.ReplaceAll(key, " ", "")
	key = strings.ReplaceAll(key, "\n", "")
	key = strings.ReplaceAll(key, "\r", "")

	// 添加 PEM 头尾，并每64个字符换行
	var formatted strings.Builder
	formatted.WriteString("-----BEGIN PUBLIC KEY-----\n")

	for i := 0; i < len(key); i += 64 {
		end := i + 64
		if end > len(key) {
			end = len(key)
		}
		formatted.WriteString(key[i:end])
		formatted.WriteString("\n")
	}

	formatted.WriteString("-----END PUBLIC KEY-----")
	return formatted.String()
}

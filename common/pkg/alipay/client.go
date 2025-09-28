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
	OutTradeNo  string  `json:"out_trade_no"`              // 商户订单号
	TotalAmount float64 `json:"total_amount"`              // 订单总金额
	Subject     string  `json:"subject"`                   // 订单标题
	Body        string  `json:"body,omitempty"`            // 订单描述
	ProductCode string  `json:"product_code"`              // 产品码
	Timeout     string  `json:"timeout_express,omitempty"` // 订单超时时间
	ReturnUrl   string  `json:"return_url,omitempty"`      // 支付成功跳转地址
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

// CreatePayment 创建支付订单
func (c *AlipayClient) CreatePayment(req *TradeCreateRequest) (*TradeCreateResponse, error) {
	bizContent, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal biz_content: %w", err)
	}

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

	// 解析响应
	var baseResp BaseResponse
	if err := json.Unmarshal(body, &baseResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查响应状态
	if baseResp.Code != "10000" {
		return nil, fmt.Errorf("alipay error: code=%s, msg=%s, sub_code=%s, sub_msg=%s",
			baseResp.Code, baseResp.Msg, baseResp.SubCode, baseResp.SubMsg)
	}

	return baseResp.Data, nil
}

// sign 对请求进行签名
func (c *AlipayClient) sign(req *BaseRequest) (string, error) {
	// 构建待签名字符串
	params := map[string]string{
		"app_id":      req.AppId,
		"method":      req.Method,
		"format":      req.Format,
		"charset":     req.Charset,
		"sign_type":   req.SignType,
		"timestamp":   req.Timestamp,
		"version":     req.Version,
		"biz_content": req.BizContent,
	}

	if req.NotifyUrl != "" {
		params["notify_url"] = req.NotifyUrl
	}

	// 排序参数
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 构建签名字符串
	var signStr strings.Builder
	for i, k := range keys {
		if i > 0 {
			signStr.WriteString("&")
		}
		signStr.WriteString(k)
		signStr.WriteString("=")
		signStr.WriteString(params[k])
	}

	// RSA2签名
	block, _ := pem.Decode([]byte(c.config.AppPrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to decode private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	hashed := sha256.Sum256([]byte(signStr.String()))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySign 验证签名
func (c *AlipayClient) VerifySign(data, sign string) error {
	// 解析支付宝公钥
	block, _ := pem.Decode([]byte(c.config.AlipayPublicKey))
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

package model

import (
	"time"

	"gorm.io/gorm"
)

// PaymentOrder 支付订单表
type PaymentOrder struct {
	ID            uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	PaymentId     string         `gorm:"column:payment_id;type:varchar(64);not null;uniqueIndex" json:"payment_id"`                 // 支付ID
	OrderId       string         `gorm:"column:order_id;type:varchar(64);not null;index" json:"order_id"`                           // 订单ID
	OutTradeNo    string         `gorm:"column:out_trade_no;type:varchar(64);not null;uniqueIndex" json:"out_trade_no"`             // 商户订单号
	UserId        uint64         `gorm:"column:user_id;type:bigint;not null;index" json:"user_id"`                                  // 用户ID
	Amount        float64        `gorm:"column:amount;type:decimal(10,2);not null" json:"amount"`                                   // 支付金额
	Subject       string         `gorm:"column:subject;type:varchar(255);not null" json:"subject"`                                  // 订单标题
	Body          string         `gorm:"column:body;type:text" json:"body"`                                                         // 订单描述
	Status        string         `gorm:"column:status;type:varchar(20);not null;default:'PENDING';index" json:"status"`             // 支付状态
	TradeNo       string         `gorm:"column:trade_no;type:varchar(64);index" json:"trade_no"`                                    // 支付宝交易号
	TradeStatus   string         `gorm:"column:trade_status;type:varchar(20)" json:"trade_status"`                                  // 支付宝交易状态
	BuyerUserId   string         `gorm:"column:buyer_user_id;type:varchar(64)" json:"buyer_user_id"`                                // 买家支付宝用户ID
	BuyerLogonId  string         `gorm:"column:buyer_logon_id;type:varchar(128)" json:"buyer_logon_id"`                             // 买家支付宝账号
	ReceiptAmount float64        `gorm:"column:receipt_amount;type:decimal(10,2)" json:"receipt_amount"`                            // 实收金额
	ProductCode   string         `gorm:"column:product_code;type:varchar(32);default:'FAST_INSTANT_TRADE_PAY'" json:"product_code"` // 产品码
	ReturnUrl     string         `gorm:"column:return_url;type:varchar(500)" json:"return_url"`                                     // 支付成功跳转地址
	NotifyUrl     string         `gorm:"column:notify_url;type:varchar(500)" json:"notify_url"`                                     // 支付结果异步通知地址
	Timeout       string         `gorm:"column:timeout;type:varchar(20)" json:"timeout"`                                            // 订单超时时间
	ClientIP      string         `gorm:"column:client_ip;type:varchar(45)" json:"client_ip"`                                        // 客户端IP
	GmtPayment    *time.Time     `gorm:"column:gmt_payment" json:"gmt_payment"`                                                     // 支付时间
	GmtClose      *time.Time     `gorm:"column:gmt_close" json:"gmt_close"`                                                         // 交易关闭时间
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName 指定表名
func (PaymentOrder) TableName() string {
	return "lxt_payment_orders"
}

// PaymentRefund 支付退款表
type PaymentRefund struct {
	ID           uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	RefundId     string         `gorm:"column:refund_id;type:varchar(64);not null;uniqueIndex" json:"refund_id"`           // 退款ID
	PaymentId    string         `gorm:"column:payment_id;type:varchar(64);not null;index" json:"payment_id"`               // 支付ID
	OrderId      string         `gorm:"column:order_id;type:varchar(64);not null;index" json:"order_id"`                   // 订单ID
	OutTradeNo   string         `gorm:"column:out_trade_no;type:varchar(64);not null;index" json:"out_trade_no"`           // 商户订单号
	OutRequestNo string         `gorm:"column:out_request_no;type:varchar(64);not null;uniqueIndex" json:"out_request_no"` // 退款单号
	UserId       uint64         `gorm:"column:user_id;type:bigint;not null;index" json:"user_id"`                          // 用户ID
	RefundAmount float64        `gorm:"column:refund_amount;type:decimal(10,2);not null" json:"refund_amount"`             // 退款金额
	RefundFee    float64        `gorm:"column:refund_fee;type:decimal(10,2)" json:"refund_fee"`                            // 退款手续费
	RefundReason string         `gorm:"column:refund_reason;type:varchar(255)" json:"refund_reason"`                       // 退款原因
	Status       string         `gorm:"column:status;type:varchar(20);not null;default:'PENDING';index" json:"status"`     // 退款状态
	RefundStatus string         `gorm:"column:refund_status;type:varchar(20)" json:"refund_status"`                        // 支付宝退款状态
	GmtRefund    *time.Time     `gorm:"column:gmt_refund" json:"gmt_refund"`                                               // 退款时间
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName 指定表名
func (PaymentRefund) TableName() string {
	return "lxt_payment_refunds"
}

// PaymentNotify 支付通知记录表
type PaymentNotify struct {
	ID            uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	NotifyId      string         `gorm:"column:notify_id;type:varchar(64);not null;uniqueIndex" json:"notify_id"`                       // 通知ID
	PaymentId     string         `gorm:"column:payment_id;type:varchar(64);not null;index" json:"payment_id"`                           // 支付ID
	NotifyType    string         `gorm:"column:notify_type;type:varchar(20);not null;index" json:"notify_type"`                         // 通知类型：PAYMENT/REFUND
	NotifyData    string         `gorm:"column:notify_data;type:longtext;not null" json:"notify_data"`                                  // 通知数据
	Sign          string         `gorm:"column:sign;type:varchar(512)" json:"sign"`                                                     // 签名
	SignType      string         `gorm:"column:sign_type;type:varchar(20)" json:"sign_type"`                                            // 签名类型
	VerifyStatus  string         `gorm:"column:verify_status;type:varchar(20);not null;default:'PENDING';index" json:"verify_status"`   // 验证状态
	ProcessStatus string         `gorm:"column:process_status;type:varchar(20);not null;default:'PENDING';index" json:"process_status"` // 处理状态
	ClientIP      string         `gorm:"column:client_ip;type:varchar(45)" json:"client_ip"`                                            // 客户端IP
	ErrorMessage  string         `gorm:"column:error_message;type:text" json:"error_message"`                                           // 错误信息
	ProcessedAt   *time.Time     `gorm:"column:processed_at" json:"processed_at"`                                                       // 处理时间
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName 指定表名
func (PaymentNotify) TableName() string {
	return "lxt_payment_notifies"
}

// PaymentConfig 支付配置表
type PaymentConfig struct {
	ID              uint64         `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	AppId           string         `gorm:"column:app_id;type:varchar(32);not null;uniqueIndex" json:"app_id"`      // 应用ID
	AppName         string         `gorm:"column:app_name;type:varchar(100);not null" json:"app_name"`             // 应用名称
	AppPrivateKey   string         `gorm:"column:app_private_key;type:text;not null" json:"app_private_key"`       // 应用私钥
	AlipayPublicKey string         `gorm:"column:alipay_public_key;type:text;not null" json:"alipay_public_key"`   // 支付宝公钥
	GatewayUrl      string         `gorm:"column:gateway_url;type:varchar(255);not null" json:"gateway_url"`       // 支付宝网关地址
	IsProd          bool           `gorm:"column:is_prod;type:tinyint(1);not null;default:0" json:"is_prod"`       // 是否生产环境
	IsEnabled       bool           `gorm:"column:is_enabled;type:tinyint(1);not null;default:1" json:"is_enabled"` // 是否启用
	NotifyUrl       string         `gorm:"column:notify_url;type:varchar(500)" json:"notify_url"`                  // 默认通知地址
	ReturnUrl       string         `gorm:"column:return_url;type:varchar(500)" json:"return_url"`                  // 默认返回地址
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName 指定表名
func (PaymentConfig) TableName() string {
	return "lxt_payment_configs"
}

// 支付状态常量
const (
	PaymentStatusPending         = "PENDING"          // 待支付
	PaymentStatusPaid            = "PAID"             // 已支付
	PaymentStatusClosed          = "CLOSED"           // 已关闭
	PaymentStatusCancelled       = "CANCELLED"        // 已取消
	PaymentStatusRefunded        = "REFUNDED"         // 已退款
	PaymentStatusPartialRefunded = "PARTIAL_REFUNDED" // 部分退款
)

// 支付宝交易状态常量
const (
	TradeStatusWaitBuyerPay = "WAIT_BUYER_PAY" // 交易创建，等待买家付款
	TradeStatusClosed       = "TRADE_CLOSED"   // 未付款交易超时关闭，或支付完成后全额退款
	TradeStatusSuccess      = "TRADE_SUCCESS"  // 交易支付成功
	TradeStatusFinished     = "TRADE_FINISHED" // 交易结束，不可退款
)

// 退款状态常量
const (
	RefundStatusPending = "PENDING" // 待退款
	RefundStatusSuccess = "SUCCESS" // 退款成功
	RefundStatusFailed  = "FAILED"  // 退款失败
	RefundStatusClosed  = "CLOSED"  // 退款关闭
)

// 通知类型常量
const (
	NotifyTypePayment = "PAYMENT" // 支付通知
	NotifyTypeRefund  = "REFUND"  // 退款通知
)

// 验证状态常量
const (
	VerifyStatusPending = "PENDING" // 待验证
	VerifyStatusSuccess = "SUCCESS" // 验证成功
	VerifyStatusFailed  = "FAILED"  // 验证失败
)

// 处理状态常量
const (
	ProcessStatusPending = "PENDING" // 待处理
	ProcessStatusSuccess = "SUCCESS" // 处理成功
	ProcessStatusFailed  = "FAILED"  // 处理失败
)

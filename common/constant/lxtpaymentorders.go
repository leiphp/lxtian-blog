package constant

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

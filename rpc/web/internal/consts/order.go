package consts

// 支付状态常量
const (
	PaymentStatusPending         = "PENDING"          // 待支付
	PaymentStatusPaid            = "PAID"             // 已支付
	PaymentStatusClosed          = "CLOSED"           // 已关闭
	PaymentStatusCancelled       = "CANCELLED"        // 已取消
	PaymentStatusRefunded        = "REFUNDED"         // 已退款
	PaymentStatusPartialRefunded = "PARTIAL_REFUNDED" // 部分退款
	PaymentStatusFAILED          = "FAILED"           // 验证失败
)

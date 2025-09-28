package payment

import (
	"context"
	"lxtian-blog/common/pkg/model"
	"lxtian-blog/common/repository"
	"time"

	"gorm.io/gorm"
)

// PaymentOrderRepository PaymentOrder表仓储接口
type PaymentOrderRepository interface {
	repository.BaseRepository[model.PaymentOrder]

	// 支付订单特有方法
	GetByPaymentId(ctx context.Context, paymentId string) (*model.PaymentOrder, error)
	GetByOrderId(ctx context.Context, orderId string) (*model.PaymentOrder, error)
	GetByOutTradeNo(ctx context.Context, outTradeNo string) (*model.PaymentOrder, error)
	GetByUserId(ctx context.Context, userId uint64, page, pageSize int) ([]*model.PaymentOrder, int64, error)
	GetByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.PaymentOrder, int64, error)
	GetByTradeNo(ctx context.Context, tradeNo string) (*model.PaymentOrder, error)

	// 更新方法
	UpdateStatus(ctx context.Context, paymentId string, status string) error
	UpdateTradeInfo(ctx context.Context, paymentId string, tradeNo, tradeStatus, buyerUserId, buyerLogonId string, receiptAmount float64, gmtPayment interface{}) error
	UpdateNotifyInfo(ctx context.Context, paymentId string, notifyData string) error

	// 统计方法
	GetCountByUserId(ctx context.Context, userId uint64) (int64, error)
	GetCountByStatus(ctx context.Context, status string) (int64, error)
	GetTotalAmountByUserId(ctx context.Context, userId uint64) (float64, error)
	GetTotalAmountByStatus(ctx context.Context, status string) (float64, error)
	GetTotalAmountByTimeRange(ctx context.Context, startTime, endTime time.Time) (float64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, paymentIds []string, status string) error
	GetExpiredOrders(ctx context.Context) ([]*model.PaymentOrder, error)
	GetOrdersByTimeRange(ctx context.Context, startTime, endTime time.Time, page, pageSize int) ([]*model.PaymentOrder, int64, error)
}

// paymentOrderRepository PaymentOrder表仓储实现
type paymentOrderRepository struct {
	*repository.TransactionalBaseRepository[model.PaymentOrder]
}

// NewPaymentOrderRepository 创建PaymentOrder仓储
func NewPaymentOrderRepository(db *gorm.DB) PaymentOrderRepository {
	return &paymentOrderRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.PaymentOrder](db),
	}
}

// GetByPaymentId 根据支付ID获取订单
func (r *paymentOrderRepository) GetByPaymentId(ctx context.Context, paymentId string) (*model.PaymentOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"payment_id": paymentId,
	})
}

// GetByOrderId 根据订单ID获取支付订单
func (r *paymentOrderRepository) GetByOrderId(ctx context.Context, orderId string) (*model.PaymentOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"order_id": orderId,
	})
}

// GetByOutTradeNo 根据商户订单号获取支付订单
func (r *paymentOrderRepository) GetByOutTradeNo(ctx context.Context, outTradeNo string) (*model.PaymentOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"out_trade_no": outTradeNo,
	})
}

// GetByUserId 根据用户ID获取支付订单列表
func (r *paymentOrderRepository) GetByUserId(ctx context.Context, userId uint64, page, pageSize int) ([]*model.PaymentOrder, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"user_id": userId,
	}, page, pageSize)
}

// GetByStatus 根据状态获取支付订单列表
func (r *paymentOrderRepository) GetByStatus(ctx context.Context, status string, page, pageSize int) ([]*model.PaymentOrder, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, page, pageSize)
}

// GetByTradeNo 根据支付宝交易号获取支付订单
func (r *paymentOrderRepository) GetByTradeNo(ctx context.Context, tradeNo string) (*model.PaymentOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"trade_no": tradeNo,
	})
}

// UpdateStatus 更新支付状态
func (r *paymentOrderRepository) UpdateStatus(ctx context.Context, paymentId string, status string) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"payment_id": paymentId},
		map[string]interface{}{"status": status},
	)
}

// UpdateTradeInfo 更新交易信息
func (r *paymentOrderRepository) UpdateTradeInfo(ctx context.Context, paymentId string, tradeNo, tradeStatus, buyerUserId, buyerLogonId string, receiptAmount float64, gmtPayment interface{}) error {
	updates := map[string]interface{}{
		"trade_no":       tradeNo,
		"trade_status":   tradeStatus,
		"buyer_user_id":  buyerUserId,
		"buyer_logon_id": buyerLogonId,
		"receipt_amount": receiptAmount,
		"status":         model.PaymentStatusPaid,
	}

	if gmtPayment != nil {
		updates["gmt_payment"] = gmtPayment
	}

	return r.UpdateByCondition(ctx,
		map[string]interface{}{"payment_id": paymentId},
		updates,
	)
}

// UpdateNotifyInfo 更新通知信息
func (r *paymentOrderRepository) UpdateNotifyInfo(ctx context.Context, paymentId string, notifyData string) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"payment_id": paymentId},
		map[string]interface{}{"notify_data": notifyData},
	)
}

// GetCountByUserId 根据用户ID统计订单数量
func (r *paymentOrderRepository) GetCountByUserId(ctx context.Context, userId uint64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"user_id": userId,
	})
}

// GetCountByStatus 根据状态统计订单数量
func (r *paymentOrderRepository) GetCountByStatus(ctx context.Context, status string) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"status": status,
	})
}

// GetTotalAmountByUserId 根据用户ID统计总金额
func (r *paymentOrderRepository) GetTotalAmountByUserId(ctx context.Context, userId uint64) (float64, error) {
	db := r.GetDB(ctx)
	var total float64

	err := db.Model(&model.PaymentOrder{}).
		Where("user_id = ? AND status = ?", userId, model.PaymentStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	return total, err
}

// GetTotalAmountByStatus 根据状态统计总金额
func (r *paymentOrderRepository) GetTotalAmountByStatus(ctx context.Context, status string) (float64, error) {
	db := r.GetDB(ctx)
	var total float64

	err := db.Model(&model.PaymentOrder{}).
		Where("status = ?", status).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	return total, err
}

// GetTotalAmountByTimeRange 根据时间范围统计总金额
func (r *paymentOrderRepository) GetTotalAmountByTimeRange(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	db := r.GetDB(ctx)
	var total float64

	err := db.Model(&model.PaymentOrder{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startTime, endTime, model.PaymentStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	return total, err
}

// BatchUpdateStatus 批量更新状态
func (r *paymentOrderRepository) BatchUpdateStatus(ctx context.Context, paymentIds []string, status string) error {
	db := r.GetDB(ctx)
	return db.Model(&model.PaymentOrder{}).
		Where("payment_id IN ?", paymentIds).
		Update("status", status).Error
}

// GetExpiredOrders 获取过期订单
func (r *paymentOrderRepository) GetExpiredOrders(ctx context.Context) ([]*model.PaymentOrder, error) {
	db := r.GetDB(ctx)
	var orders []*model.PaymentOrder

	// 查询创建时间超过30分钟且状态为待支付的订单
	cutoffTime := time.Now().Add(-30 * time.Minute)
	err := db.Where("status = ? AND created_at < ?", model.PaymentStatusPending, cutoffTime).
		Find(&orders).Error

	return orders, err
}

// GetOrdersByTimeRange 根据时间范围获取订单
func (r *paymentOrderRepository) GetOrdersByTimeRange(ctx context.Context, startTime, endTime time.Time, page, pageSize int) ([]*model.PaymentOrder, int64, error) {
	db := r.GetDB(ctx)
	var orders []*model.PaymentOrder
	var total int64

	query := db.Where("created_at BETWEEN ? AND ?", startTime, endTime)

	// 获取总数
	if err := query.Model(&model.PaymentOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

package payment_repo

import (
	"context"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository"
	"time"

	"gorm.io/gorm"
)

// PaymentOrderRepository PaymentOrder表仓储接口
type PaymentOrderRepository interface {
	repository.BaseRepository[model.LxtPaymentOrder]

	// 支付订单特有方法
	GetByPaymentId(ctx context.Context, paymentId string) (*model.LxtPaymentOrder, error)
	GetByOrderSn(ctx context.Context, orderSn string) (*model.LxtPaymentOrder, error)
	GetByOutTradeNo(ctx context.Context, outTradeNo string) (*model.LxtPaymentOrder, error)
	GetByUserId(ctx context.Context, userId uint64, page, pageSize int, keywords string) ([]*model.LxtPaymentOrder, int64, error)
	GetByStatus(ctx context.Context, status string, page, pageSize int, keywords string) ([]*model.LxtPaymentOrder, int64, error)
	GetByTradeNo(ctx context.Context, tradeNo string) (*model.LxtPaymentOrder, error)

	// 更新方法
	UpdateStatus(ctx context.Context, paymentId string, status string) error
	UpdateTradeInfo(ctx context.Context, paymentId string, tradeNo, tradeStatus, buyerUserId, buyerLogonId string, receiptAmount float64, gmtPayment interface{}) error
	UpdateNotifyInfo(ctx context.Context, paymentId string, notifyData string) error

	// 删除方法
	SoftDeleteByOrderSn(ctx context.Context, orderSn string) error

	// 统计方法
	GetCountByUserId(ctx context.Context, userId uint64) (int64, error)
	GetCountByStatus(ctx context.Context, status string) (int64, error)
	GetTotalAmountByUserId(ctx context.Context, userId uint64) (float64, error)
	GetTotalAmountByStatus(ctx context.Context, status string) (float64, error)
	GetTotalAmountByTimeRange(ctx context.Context, startTime, endTime time.Time) (float64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, paymentIds []string, status string) error
	GetExpiredOrders(ctx context.Context) ([]*model.LxtPaymentOrder, error)
	GetOrdersByTimeRange(ctx context.Context, startTime, endTime time.Time, page, pageSize int) ([]*model.LxtPaymentOrder, int64, error)

	// 支付通知相关方法
	FindPaymentOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*model.LxtPaymentOrder, error)
	FindPaymentNotifyByNotifyId(ctx context.Context, notifyId string) (*model.LxtPaymentNotify, error)
	UpdatePaymentNotify(ctx context.Context, notify *model.LxtPaymentNotify) error
	UpdatePaymentNotifyVerifyStatus(ctx context.Context, notifyId string, verifyStatus string) error
	UpdatePaymentNotifyProcessStatus(ctx context.Context, notifyId string, processStatus string, errorMsg string) error
	UpdatePaymentOrderStatus(ctx context.Context, paymentId string, status string) error
	UpdatePaymentOrderTradeInfo(ctx context.Context, paymentId string, tradeNo, tradeStatus, buyerUserId, buyerLogonId string, receiptAmount float64, gmtPayment *time.Time) error
}

// paymentOrderRepository PaymentOrder表仓储实现
type paymentOrderRepository struct {
	*repository.TransactionalBaseRepository[model.LxtPaymentOrder]
}

// NewPaymentOrderRepository 创建PaymentOrder仓储
func NewPaymentOrderRepository(db *gorm.DB) PaymentOrderRepository {
	return &paymentOrderRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.LxtPaymentOrder](db),
	}
}

// GetByPaymentId 根据支付ID获取订单
func (r *paymentOrderRepository) GetByPaymentId(ctx context.Context, paymentId string) (*model.LxtPaymentOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"payment_id": paymentId,
	})
}

// GetByOrderSn 根据订单Sn获取支付订单
func (r *paymentOrderRepository) GetByOrderSn(ctx context.Context, orderSn string) (*model.LxtPaymentOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"order_sn": orderSn,
	})
}

// GetByOutTradeNo 根据商户订单号获取支付订单
func (r *paymentOrderRepository) GetByOutTradeNo(ctx context.Context, outTradeNo string) (*model.LxtPaymentOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"out_trade_no": outTradeNo,
	})
}

// GetByUserId 根据用户ID获取支付订单列表
func (r *paymentOrderRepository) GetByUserId(ctx context.Context, userId uint64, page, pageSize int, keywords string) ([]*model.LxtPaymentOrder, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"user_id": userId,
	}, page, pageSize, "", "")
}

// GetByStatus 根据状态获取支付订单列表
func (r *paymentOrderRepository) GetByStatus(ctx context.Context, status string, page, pageSize int, keywords string) ([]*model.LxtPaymentOrder, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, page, pageSize, "", "")
}

// GetByTradeNo 根据支付宝交易号获取支付订单
func (r *paymentOrderRepository) GetByTradeNo(ctx context.Context, tradeNo string) (*model.LxtPaymentOrder, error) {
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
		"status":         constant.PaymentStatusPaid,
	}

	if gmtPayment != nil {
		updates["pay_time"] = gmtPayment
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

	err := db.Model(&model.LxtPaymentOrder{}).
		Where("user_id = ? AND status = ?", userId, constant.PaymentStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	return total, err
}

// GetTotalAmountByStatus 根据状态统计总金额
func (r *paymentOrderRepository) GetTotalAmountByStatus(ctx context.Context, status string) (float64, error) {
	db := r.GetDB(ctx)
	var total float64

	err := db.Model(&model.LxtPaymentOrder{}).
		Where("status = ?", status).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	return total, err
}

// GetTotalAmountByTimeRange 根据时间范围统计总金额
func (r *paymentOrderRepository) GetTotalAmountByTimeRange(ctx context.Context, startTime, endTime time.Time) (float64, error) {
	db := r.GetDB(ctx)
	var total float64

	err := db.Model(&model.LxtPaymentOrder{}).
		Where("created_at BETWEEN ? AND ? AND status = ?", startTime, endTime, constant.PaymentStatusPaid).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	return total, err
}

// BatchUpdateStatus 批量更新状态
func (r *paymentOrderRepository) BatchUpdateStatus(ctx context.Context, paymentIds []string, status string) error {
	db := r.GetDB(ctx)
	return db.Model(&model.LxtPaymentOrder{}).
		Where("payment_id IN ?", paymentIds).
		Update("status", status).Error
}

// GetExpiredOrders 获取过期订单
func (r *paymentOrderRepository) GetExpiredOrders(ctx context.Context) ([]*model.LxtPaymentOrder, error) {
	db := r.GetDB(ctx)
	var orders []*model.LxtPaymentOrder

	// 查询创建时间超过30分钟且状态为待支付的订单
	cutoffTime := time.Now().Add(-30 * time.Minute)
	err := db.Where("status = ? AND created_at < ?", constant.PaymentStatusPending, cutoffTime).
		Find(&orders).Error

	return orders, err
}

// GetOrdersByTimeRange 根据时间范围获取订单
func (r *paymentOrderRepository) GetOrdersByTimeRange(ctx context.Context, startTime, endTime time.Time, page, pageSize int) ([]*model.LxtPaymentOrder, int64, error) {
	db := r.GetDB(ctx)
	var orders []*model.LxtPaymentOrder
	var total int64

	query := db.Where("created_at BETWEEN ? AND ?", startTime, endTime)

	// 获取总数
	if err := query.Model(&model.LxtPaymentOrder{}).Count(&total).Error; err != nil {
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

// FindPaymentOrderByOutTradeNo 根据商户订单号查找支付订单
func (r *paymentOrderRepository) FindPaymentOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*model.LxtPaymentOrder, error) {
	return r.GetByOutTradeNo(ctx, outTradeNo)
}

// FindPaymentNotifyByNotifyId 根据通知ID查找支付通知
func (r *paymentOrderRepository) FindPaymentNotifyByNotifyId(ctx context.Context, notifyId string) (*model.LxtPaymentNotify, error) {
	db := r.GetDB(ctx)
	var notify model.LxtPaymentNotify

	err := db.Where("notify_id = ?", notifyId).First(&notify).Error
	if err != nil {
		return nil, err
	}

	return &notify, nil
}

// UpdatePaymentNotify 更新支付通知记录
func (r *paymentOrderRepository) UpdatePaymentNotify(ctx context.Context, notify *model.LxtPaymentNotify) error {
	db := r.GetDB(ctx)
	return db.Model(&model.LxtPaymentNotify{}).
		Where("notify_id = ?", notify.NotifyID).
		Updates(notify).Error
}

// UpdatePaymentNotifyVerifyStatus 更新支付通知验证状态
func (r *paymentOrderRepository) UpdatePaymentNotifyVerifyStatus(ctx context.Context, notifyId string, verifyStatus string) error {
	db := r.GetDB(ctx)
	return db.Model(&model.LxtPaymentNotify{}).
		Where("notify_id = ?", notifyId).
		Update("verify_status", verifyStatus).Error
}

// UpdatePaymentNotifyProcessStatus 更新支付通知处理状态
func (r *paymentOrderRepository) UpdatePaymentNotifyProcessStatus(ctx context.Context, notifyId string, processStatus string, errorMsg string) error {
	db := r.GetDB(ctx)
	updates := map[string]interface{}{
		"process_status": processStatus,
	}

	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	if processStatus == constant.ProcessStatusSuccess || processStatus == constant.ProcessStatusFailed {
		now := time.Now()
		updates["processed_at"] = &now
	}

	return db.Model(&model.LxtPaymentNotify{}).
		Where("notify_id = ?", notifyId).
		Updates(updates).Error
}

// UpdatePaymentOrderStatus 更新支付订单状态
func (r *paymentOrderRepository) UpdatePaymentOrderStatus(ctx context.Context, paymentId string, status string) error {
	return r.UpdateStatus(ctx, paymentId, status)
}

// UpdatePaymentOrderTradeInfo 更新支付订单交易信息
func (r *paymentOrderRepository) UpdatePaymentOrderTradeInfo(ctx context.Context, paymentId string, tradeNo, tradeStatus, buyerUserId, buyerLogonId string, receiptAmount float64, gmtPayment *time.Time) error {
	return r.UpdateTradeInfo(ctx, paymentId, tradeNo, tradeStatus, buyerUserId, buyerLogonId, receiptAmount, gmtPayment)
}

// SoftDeleteByOrderSn 根据订单编号软删除订单（设置deleted_at不为空）
func (r *paymentOrderRepository) SoftDeleteByOrderSn(ctx context.Context, orderSn string) error {
	db := r.GetDB(ctx)
	return db.Where("order_sn = ?", orderSn).Delete(&model.LxtPaymentOrder{}).Error
}

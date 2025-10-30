package payment_repo

import (
	"context"
	"database/sql"

	"lxtian-blog/common/model"
	"lxtian-blog/common/repository"

	"gorm.io/gorm"
)

type LxtPaymentRefundsRepo interface {
	repository.BaseRepository[model.LxtPaymentRefund]

	// 退款相关方法
	FindPaymentOrderByPaymentId(ctx context.Context, paymentId string) (*model.LxtPaymentOrder, error)
	FindPaymentOrderByOrderId(ctx context.Context, orderId string) (*model.LxtPaymentOrder, error)
	FindPaymentOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*model.LxtPaymentOrder, error)
	FindPaymentRefundByOutRequestNo(ctx context.Context, outRequestNo string) (*model.LxtPaymentRefund, error)
	InsertPaymentRefund(ctx context.Context, refund *model.LxtPaymentRefund) (sql.Result, error)
	UpdatePaymentRefund(ctx context.Context, refund *model.LxtPaymentRefund) error
	UpdatePaymentRefundStatus(ctx context.Context, refundId string, status string) error
	UpdatePaymentOrderStatus(ctx context.Context, paymentId string, status string) error
}

type lxtPaymentRefundsRepo struct {
	*repository.TransactionalBaseRepository[model.LxtPaymentRefund]
}

func NewLxtPaymentRefundsRepo(db *gorm.DB) LxtPaymentRefundsRepo {
	return &lxtPaymentRefundsRepo{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.LxtPaymentRefund](db),
	}
}

// FindPaymentOrderByPaymentId 根据支付ID查找支付订单
func (r *lxtPaymentRefundsRepo) FindPaymentOrderByPaymentId(ctx context.Context, paymentId string) (*model.LxtPaymentOrder, error) {
	db := r.GetDB(ctx)
	var order model.LxtPaymentOrder

	err := db.Where("payment_id = ?", paymentId).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// FindPaymentOrderByOrderId 根据订单ID查找支付订单
func (r *lxtPaymentRefundsRepo) FindPaymentOrderByOrderId(ctx context.Context, orderId string) (*model.LxtPaymentOrder, error) {
	db := r.GetDB(ctx)
	var order model.LxtPaymentOrder

	err := db.Where("order_id = ?", orderId).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// FindPaymentOrderByOutTradeNo 根据商户订单号查找支付订单
func (r *lxtPaymentRefundsRepo) FindPaymentOrderByOutTradeNo(ctx context.Context, outTradeNo string) (*model.LxtPaymentOrder, error) {
	db := r.GetDB(ctx)
	var order model.LxtPaymentOrder

	err := db.Where("out_trade_no = ?", outTradeNo).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// FindPaymentRefundByOutRequestNo 根据退款单号查找退款记录
func (r *lxtPaymentRefundsRepo) FindPaymentRefundByOutRequestNo(ctx context.Context, outRequestNo string) (*model.LxtPaymentRefund, error) {
	db := r.GetDB(ctx)
	var refund model.LxtPaymentRefund

	err := db.Where("out_request_no = ?", outRequestNo).First(&refund).Error
	if err != nil {
		return nil, err
	}

	return &refund, nil
}

// InsertPaymentRefund 插入退款记录
func (r *lxtPaymentRefundsRepo) InsertPaymentRefund(ctx context.Context, refund *model.LxtPaymentRefund) (sql.Result, error) {
	db := r.GetDB(ctx)
	result := db.Create(refund)
	return nil, result.Error
}

// UpdatePaymentRefund 更新退款记录
func (r *lxtPaymentRefundsRepo) UpdatePaymentRefund(ctx context.Context, refund *model.LxtPaymentRefund) error {
	db := r.GetDB(ctx)
	return db.Model(&model.LxtPaymentRefund{}).
		Where("refund_id = ?", refund.RefundID).
		Updates(refund).Error
}

// UpdatePaymentRefundStatus 更新退款状态
func (r *lxtPaymentRefundsRepo) UpdatePaymentRefundStatus(ctx context.Context, refundId string, status string) error {
	db := r.GetDB(ctx)
	return db.Model(&model.LxtPaymentRefund{}).
		Where("refund_id = ?", refundId).
		Update("status", status).Error
}

// UpdatePaymentOrderStatus 更新支付订单状态
func (r *lxtPaymentRefundsRepo) UpdatePaymentOrderStatus(ctx context.Context, paymentId string, status string) error {
	db := r.GetDB(ctx)
	return db.Model(&model.LxtPaymentOrder{}).
		Where("payment_id = ?", paymentId).
		Update("status", status).Error
}

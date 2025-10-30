package web_repo

import (
	"context"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/repository"
	"time"

	"gorm.io/gorm"
)

// TxyOrderRepository TxyOrder表仓储接口
type TxyOrderRepository interface {
	repository.BaseRepository[mysql.TxyOrder]

	// 业务订单特有方法
	GetByOrderSn(ctx context.Context, orderSn string) (*mysql.TxyOrder, error)
	GetByOutTradeNo(ctx context.Context, outTradeNo string) (*mysql.TxyOrder, error)
	GetByUserId(ctx context.Context, userId uint64, page, pageSize int, keywords string) ([]*mysql.TxyOrder, int64, error)
	GetByPayType(ctx context.Context, payType int64, page, pageSize int, keywords string) ([]*mysql.TxyOrder, int64, error)
	GetByStatus(ctx context.Context, status int64, page, pageSize int, keywords string) ([]*mysql.TxyOrder, int64, error)
	GetOrdersByTimeRange(ctx context.Context, startTime, endTime int64, page, pageSize int) ([]*mysql.TxyOrder, int64, error)

	// 更新方法
	UpdateStatus(ctx context.Context, orderId uint64, status int64) error
	UpdatePayInfo(ctx context.Context, orderId uint64, payMoney float64, payTime int64) error
	UpdateRemark(ctx context.Context, orderId uint64, remark string) error

	// 统计方法
	GetCountByUserId(ctx context.Context, userId uint64) (int64, error)
	GetCountByPayType(ctx context.Context, payType int64) (int64, error)
	GetCountByStatus(ctx context.Context, status int64) (int64, error)
	GetTotalAmountByUserId(ctx context.Context, userId uint64) (float64, error)
	GetTotalAmountByPayType(ctx context.Context, payType int64) (float64, error)
	GetTotalAmountByTimeRange(ctx context.Context, startTime, endTime int64) (float64, error)

	// 批量操作
	BatchUpdateStatus(ctx context.Context, orderIds []uint64, status int64) error
	GetExpiredOrders(ctx context.Context, days int) ([]*mysql.TxyOrder, error)
}

// txyOrderRepository TxyOrder表仓储实现
type txyOrderRepository struct {
	*repository.TransactionalBaseRepository[mysql.TxyOrder]
}

// NewTxyOrderRepository 创建TxyOrder仓储
func NewTxyOrderRepository(db *gorm.DB) TxyOrderRepository {
	return &txyOrderRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[mysql.TxyOrder](db),
	}
}

// GetByOrderSn 根据订单号获取订单
func (r *txyOrderRepository) GetByOrderSn(ctx context.Context, orderSn string) (*mysql.TxyOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"order_sn": orderSn,
	})
}

// GetByOutTradeNo 根据商户订单号获取订单
func (r *txyOrderRepository) GetByOutTradeNo(ctx context.Context, outTradeNo string) (*mysql.TxyOrder, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"out_trade_no": outTradeNo,
	})
}

// GetByUserId 根据用户ID获取订单列表
func (r *txyOrderRepository) GetByUserId(ctx context.Context, userId uint64, page, pageSize int, keywords string) ([]*mysql.TxyOrder, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"user_id": userId,
	}, page, pageSize, "", "")
}

// GetByPayType 根据支付类型获取订单列表
func (r *txyOrderRepository) GetByPayType(ctx context.Context, payType int64, page, pageSize int, keywords string) ([]*mysql.TxyOrder, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"pay_type": payType,
	}, page, pageSize, "", "")
}

// GetByStatus 根据状态获取订单列表
func (r *txyOrderRepository) GetByStatus(ctx context.Context, status int64, page, pageSize int, keywords string) ([]*mysql.TxyOrder, int64, error) {
	return r.GetList(ctx, map[string]interface{}{
		"status": status,
	}, page, pageSize, "", "")
}

// GetOrdersByTimeRange 根据时间范围获取订单列表
func (r *txyOrderRepository) GetOrdersByTimeRange(ctx context.Context, startTime, endTime int64, page, pageSize int) ([]*mysql.TxyOrder, int64, error) {
	db := r.GetDB(ctx)
	var orders []*mysql.TxyOrder
	var total int64

	query := db.Where("ctime BETWEEN ? AND ?", startTime, endTime)

	// 获取总数
	if err := query.Model(&mysql.TxyOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	if err := query.Order("ctime DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// UpdateStatus 更新订单状态
func (r *txyOrderRepository) UpdateStatus(ctx context.Context, orderId uint64, status int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": orderId},
		map[string]interface{}{"status": status},
	)
}

// UpdatePayInfo 更新支付信息
func (r *txyOrderRepository) UpdatePayInfo(ctx context.Context, orderId uint64, payMoney float64, payTime int64) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": orderId},
		map[string]interface{}{
			"pay_money": payMoney,
			"pay_time":  payTime,
			"status":    1, // 1表示已支付
		},
	)
}

// UpdateRemark 更新备注
func (r *txyOrderRepository) UpdateRemark(ctx context.Context, orderId uint64, remark string) error {
	return r.UpdateByCondition(ctx,
		map[string]interface{}{"id": orderId},
		map[string]interface{}{"remark": remark},
	)
}

// GetCountByUserId 根据用户ID统计订单数量
func (r *txyOrderRepository) GetCountByUserId(ctx context.Context, userId uint64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"user_id": userId,
	})
}

// GetCountByPayType 根据支付类型统计订单数量
func (r *txyOrderRepository) GetCountByPayType(ctx context.Context, payType int64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"pay_type": payType,
	})
}

// GetCountByStatus 根据状态统计订单数量
func (r *txyOrderRepository) GetCountByStatus(ctx context.Context, status int64) (int64, error) {
	return r.Count(ctx, map[string]interface{}{
		"status": status,
	})
}

// GetTotalAmountByUserId 根据用户ID统计总金额
func (r *txyOrderRepository) GetTotalAmountByUserId(ctx context.Context, userId uint64) (float64, error) {
	db := r.GetDB(ctx)
	var total float64

	err := db.Model(&mysql.TxyOrder{}).
		Where("user_id = ? AND status = ?", userId, 1). // 1表示已支付
		Select("COALESCE(SUM(pay_money), 0)").
		Scan(&total).Error

	return total, err
}

// GetTotalAmountByPayType 根据支付类型统计总金额
func (r *txyOrderRepository) GetTotalAmountByPayType(ctx context.Context, payType int64) (float64, error) {
	db := r.GetDB(ctx)
	var total float64

	err := db.Model(&mysql.TxyOrder{}).
		Where("pay_type = ? AND status = ?", payType, 1). // 1表示已支付
		Select("COALESCE(SUM(pay_money), 0)").
		Scan(&total).Error

	return total, err
}

// GetTotalAmountByTimeRange 根据时间范围统计总金额
func (r *txyOrderRepository) GetTotalAmountByTimeRange(ctx context.Context, startTime, endTime int64) (float64, error) {
	db := r.GetDB(ctx)
	var total float64

	err := db.Model(&mysql.TxyOrder{}).
		Where("ctime BETWEEN ? AND ? AND status = ?", startTime, endTime, 1). // 1表示已支付
		Select("COALESCE(SUM(pay_money), 0)").
		Scan(&total).Error

	return total, err
}

// BatchUpdateStatus 批量更新状态
func (r *txyOrderRepository) BatchUpdateStatus(ctx context.Context, orderIds []uint64, status int64) error {
	db := r.GetDB(ctx)
	return db.Model(&mysql.TxyOrder{}).
		Where("id IN ?", orderIds).
		Update("status", status).Error
}

// GetExpiredOrders 获取过期订单
func (r *txyOrderRepository) GetExpiredOrders(ctx context.Context, days int) ([]*mysql.TxyOrder, error) {
	db := r.GetDB(ctx)
	var orders []*mysql.TxyOrder

	// 计算过期时间戳
	cutoffTime := time.Now().Unix() - int64(days*24*3600)

	err := db.Where("ctime < ? AND status = ?", cutoffTime, 0). // 0表示待支付状态
									Order("ctime ASC").
									Find(&orders).Error

	return orders, err
}

package payment_repo

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository"
)

// LxtPaymentGoodsRepo 商品表仓储接口
type LxtPaymentGoodsRepo interface {
	repository.BaseRepository[model.LxtPaymentGood]

	GetById(ctx context.Context, Id int64) (*model.LxtPaymentGood, error)
	GetByOrderId(ctx context.Context, orderId string) (*model.LxtPaymentGood, error)
}

// lxtPaymentGoodsRepo 商品表仓储实现
type lxtPaymentGoodsRepo struct {
	*repository.TransactionalBaseRepository[model.LxtPaymentGood]
}

// NewLxtPaymentGoodsRepo 创建商品表仓储
func NewLxtPaymentGoodsRepo(db *gorm.DB) LxtPaymentGoodsRepo {
	return &lxtPaymentGoodsRepo{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.LxtPaymentGood](db),
	}
}

// GetById 根据ID获取订单
func (r *lxtPaymentGoodsRepo) GetById(ctx context.Context, Id int64) (*model.LxtPaymentGood, error) {
	fmt.Println("id:", Id)
	return r.GetByCondition(ctx, map[string]interface{}{
		"id": Id,
	})
}

// GetByOrderId 根据订单ID获取支付订单
func (r *lxtPaymentGoodsRepo) GetByOrderId(ctx context.Context, orderId string) (*model.LxtPaymentGood, error) {
	return r.GetByCondition(ctx, map[string]interface{}{
		"order_id": orderId,
	})
}

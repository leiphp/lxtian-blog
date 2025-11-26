package weblogic

import (
	"context"
	"lxtian-blog/rpc/web/internal/consts"
	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderStatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderStatLogic {
	return &OrderStatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderStatLogic) OrderStat(in *web.OrderStatReq) (*web.OrderStatResp, error) {
	var totalAmount float64
	var count int64
	var monthAmount float64

	// 统计所有已支付订单的总金额和数量
	var result struct {
		TotalAmount float64 `gorm:"column:total_amount"`
		Count       int64   `gorm:"column:count"`
	}
	err := l.svcCtx.DB.Table("txy_order").
		Where("status = ?", consts.PaymentStatusPaid).
		Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as count").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	totalAmount = result.TotalAmount
	count = result.Count

	// 计算本月的开始和结束时间
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// 统计本月已支付订单的总金额
	err = l.svcCtx.DB.Table("txy_order").
		Where("status = ? AND created_at >= ? AND created_at <= ?", consts.PaymentStatusPaid, startOfMonth, endOfMonth).
		Select("COALESCE(SUM(amount), 0) as month_amount").
		Scan(&monthAmount).Error
	if err != nil {
		return nil, err
	}

	return &web.OrderStatResp{
		TotalAmount: float32(totalAmount),
		Count:       count,
		MonthAmount: float32(monthAmount),
	}, nil
}

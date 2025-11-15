package logic

import (
	"context"

	"lxtian-blog/common/constant"
	"lxtian-blog/common/repository/payment_repo"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"

	"github.com/shopspring/decimal"
)

type OrdersStatisticsLogic struct {
	*BaseLogic
	paymentService payment_repo.PaymentOrderRepository
}

func NewOrdersStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrdersStatisticsLogic {
	return &OrdersStatisticsLogic{
		BaseLogic:      NewBaseLogic(ctx, svcCtx),
		paymentService: payment_repo.NewPaymentOrderRepository(svcCtx.DB),
	}
}

// 支付订单统计
func (l *OrdersStatisticsLogic) OrdersStatistics(in *payment.OrdersStatisticsReq) (*payment.OrdersStatisticsResp, error) {
	var total, paid, pending int64
	var payAmount float64
	var err error

	// 如果指定了用户ID，则统计该用户的订单
	if in.UserId > 0 {
		// 获取总数
		total, err = l.paymentService.GetCountByUserId(l.ctx, uint64(in.UserId))
		if err != nil {
			l.Errorf("Failed to get total count by user_id: %v", err)
			return nil, err
		}

		// 获取已完成数量（已支付）
		condition := map[string]interface{}{
			"user_id": uint64(in.UserId),
			"status":  constant.PaymentStatusPaid,
		}
		paid, err = l.paymentService.Count(l.ctx, condition)
		if err != nil {
			l.Errorf("Failed to get finish count: %v", err)
			return nil, err
		}

		// 获取待支付数量
		pendingCondition := map[string]interface{}{
			"user_id": uint64(in.UserId),
			"status":  constant.PaymentStatusPending,
		}
		pending, err = l.paymentService.Count(l.ctx, pendingCondition)
		if err != nil {
			l.Errorf("Failed to get pending count: %v", err)
			return nil, err
		}

		// 获取支付总金额
		payAmount, err = l.paymentService.GetTotalAmountByUserId(l.ctx, uint64(in.UserId))
		if err != nil {
			l.Errorf("Failed to get total amount: %v", err)
			return nil, err
		}
	} else {
		// 如果没有指定用户ID，则统计所有订单
		// 获取总数
		total, err = l.paymentService.Count(l.ctx, map[string]interface{}{})
		if err != nil {
			l.Errorf("Failed to get total count: %v", err)
			return nil, err
		}

		// 获取已完成数量（已支付）
		paid, err = l.paymentService.GetCountByStatus(l.ctx, constant.PaymentStatusPaid)
		if err != nil {
			l.Errorf("Failed to get finish count: %v", err)
			return nil, err
		}

		// 获取待支付数量
		pending, err = l.paymentService.GetCountByStatus(l.ctx, constant.PaymentStatusPending)
		if err != nil {
			l.Errorf("Failed to get pending count: %v", err)
			return nil, err
		}

		// 获取支付总金额
		payAmount, err = l.paymentService.GetTotalAmountByStatus(l.ctx, constant.PaymentStatusPaid)
		if err != nil {
			l.Errorf("Failed to get total amount: %v", err)
			return nil, err
		}
	}

	// 使用 decimal 库进行精确计算，避免浮点数精度问题
	// 将 float64 转换为 decimal，保留2位小数，然后转换回 float32
	payAmountDecimal := decimal.NewFromFloat(payAmount).Round(2)
	payAmountFloat32, _ := payAmountDecimal.Float64()

	return &payment.OrdersStatisticsResp{
		Total:     total,
		Paid:      paid,
		Pending:   pending,
		PayAmount: float32(payAmountFloat32),
	}, nil
}

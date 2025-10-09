package logic

import (
	"context"

	"lxtian-blog/common/constant"
	paymentSvc "lxtian-blog/common/repository/payment"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type OrdersStatisticsLogic struct {
	*BaseLogic
	paymentService paymentSvc.PaymentOrderRepository
}

func NewOrdersStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrdersStatisticsLogic {
	return &OrdersStatisticsLogic{
		BaseLogic:      NewBaseLogic(ctx, svcCtx),
		paymentService: paymentSvc.NewPaymentOrderRepository(svcCtx.DB),
	}
}

// 支付订单统计
func (l *OrdersStatisticsLogic) OrdersStatistics(in *payment.OrdersStatisticsReq) (*payment.OrdersStatisticsResp, error) {
	var total, finish, pending int64
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
		finish, err = l.paymentService.Count(l.ctx, condition)
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
		finish, err = l.paymentService.GetCountByStatus(l.ctx, constant.PaymentStatusPaid)
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

	return &payment.OrdersStatisticsResp{
		Total:     total,
		Finish:    finish,
		Pending:   pending,
		PayAmount: float32(payAmount),
	}, nil
}

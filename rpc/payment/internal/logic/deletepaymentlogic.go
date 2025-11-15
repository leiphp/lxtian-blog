package logic

import (
	"context"
	"fmt"
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository/payment_repo"

	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePaymentLogic {
	return &DeletePaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除支付订单
func (l *DeletePaymentLogic) DeletePayment(in *payment.DeletePaymentReq) (*payment.DeletePaymentResp, error) {
	// 参数验证
	if in.OrderSn == "" {
		return &payment.DeletePaymentResp{
			Success: false,
			Message: "订单SN不能为空",
		}, fmt.Errorf("order_sn is required")
	}

	var paymentOrder *model.LxtPaymentOrder
	var err error
	paymentService := payment_repo.NewPaymentOrderRepository(l.svcCtx.DB)
	paymentOrder, err = paymentService.GetByOrderSn(l.ctx, in.OrderSn)
	if err != nil {
		l.Errorf("Failed to find payment order: %v", err)
		return &payment.DeletePaymentResp{
			Success: false,
			Message: "支付订单SN不存在",
		}, fmt.Errorf("payment order not found: %w", err)
	}

	// 验证订单是否属于该用户
	if paymentOrder.UserID != int64(in.UserId) {
		l.Errorf("User %d attempted to delete order %s belonging to user %d", in.UserId, in.OrderSn, paymentOrder.UserID)
		return &payment.DeletePaymentResp{
			Success: false,
			Message: "无权删除该订单",
		}, fmt.Errorf("order does not belong to user")
	}

	// 更新本地订单状态
	err = paymentService.SoftDeleteByOrderSn(l.ctx, in.OrderSn)
	if err != nil {
		l.Errorf("Failed to SoftDeleteByOrderSn: %v", err)
		// 即使本地更新失败，支付宝那边已经取消了，所以仍然返回成功
	}

	// 记录日志
	l.Infof("Deleted payment order: paymentId=%s, orderSn=%s",
		paymentOrder.PaymentID, paymentOrder.OrderSn)

	return &payment.DeletePaymentResp{
		Success: true,
		Message: "删除成功",
	}, nil
}

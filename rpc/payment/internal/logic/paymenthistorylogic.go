package logic

import (
	"context"
	"encoding/json"
	"fmt"
	paymentSvc "lxtian-blog/common/repository/payment"
	"time"

	"lxtian-blog/common/model"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type PaymentHistoryLogic struct {
	*BaseLogic
	paymentService paymentSvc.PaymentOrderRepository
}

func NewPaymentHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentHistoryLogic {
	return &PaymentHistoryLogic{
		BaseLogic:      NewBaseLogic(ctx, svcCtx),
		paymentService: paymentSvc.NewPaymentOrderRepository(svcCtx.DB),
	}
}

func (l *PaymentHistoryLogic) PaymentHistory(in *payment.PaymentHistoryReq) (*payment.PaymentHistoryResp, error) {
	// 参数验证
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	if in.PageSize > 100 {
		in.PageSize = 100 // 限制最大每页数量
	}

	// 构建查询条件
	condition := make(map[string]interface{})

	if in.UserId > 0 {
		condition["user_id"] = in.UserId
	}

	if in.Status != "" {
		condition["status"] = in.Status
	}

	// 使用基础仓储的 GetList 方法（已修复总数查询问题）
	paymentOrders, total, err := l.paymentService.GetList(l.ctx, condition, int(in.Page), int(in.PageSize), in.Keywords)
	if err != nil {
		l.Errorf("Failed to get payment orders: %v", err)
		return nil, fmt.Errorf("failed to get payment orders: %w", err)
	}

	// 构建响应数据
	var listData []map[string]interface{}
	for _, order := range paymentOrders {
		item := l.buildPaymentOrderItem(order)
		listData = append(listData, item)
	}

	// 转换为JSON字符串
	listJson, err := json.Marshal(listData)
	if err != nil {
		l.Errorf("Failed to marshal payment orders: %v", err)
		return nil, fmt.Errorf("failed to marshal payment orders: %w", err)
	}

	return &payment.PaymentHistoryResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint64(total), // 使用数据库查询的真实总数
		List:     string(listJson),
	}, nil
}

// 按时间范围过滤
func (l *PaymentHistoryLogic) filterByTimeRange(orders []*model.LxtPaymentOrders, startTime, endTime string) []*model.LxtPaymentOrders {
	var filtered []*model.LxtPaymentOrders

	var start, end time.Time
	var err error

	if startTime != "" {
		start, err = time.Parse("2006-01-02", startTime)
		if err != nil {
			l.Errorf("Failed to parse start_time: %v", err)
			return orders // 解析失败时返回原数据
		}
	}

	if endTime != "" {
		end, err = time.Parse("2006-01-02", endTime)
		if err != nil {
			l.Errorf("Failed to parse end_time: %v", err)
			return orders // 解析失败时返回原数据
		}
		// 设置为当天结束时间
		end = end.Add(24*time.Hour - time.Second)
	}

	for _, order := range orders {
		// 检查创建时间是否在范围内
		if !start.IsZero() && order.CreatedAt.Before(start) {
			continue
		}
		if !end.IsZero() && order.CreatedAt.After(end) {
			continue
		}
		filtered = append(filtered, order)
	}

	return filtered
}

// 构建支付订单项
func (l *PaymentHistoryLogic) buildPaymentOrderItem(order *model.LxtPaymentOrders) map[string]interface{} {
	item := map[string]interface{}{
		"id":           order.Id,
		"payment_id":   order.PaymentId,
		"order_sn":     order.OrderSn,
		"out_trade_no": order.OutTradeNo,
		"user_id":      order.UserId,
		"amount":       order.Amount,
		"subject":      order.Subject,
		"body":         order.Body,
		"status":       order.Status,
		"trade_no":     order.TradeNo,
		"trade_status": order.TradeStatus,
		"product_code": order.ProductCode,
		"return_url":   order.ReturnUrl,
		"notify_url":   order.NotifyUrl,
		"timeout":      order.Timeout,
		"created_at":   order.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if order.BuyerUserId != "" {
		item["buyer_user_id"] = order.BuyerUserId
	}
	if order.BuyerLogonId != "" {
		item["buyer_logon_id"] = order.BuyerLogonId
	}
	if order.ReceiptAmount != 0 {
		item["receipt_amount"] = order.ReceiptAmount
	}
	if order.PayTime.Valid {
		item["pay_time"] = order.PayTime.Time.Format("2006-01-02 15:04:05")
	}
	if order.CloseTime.Valid {
		item["close_time"] = order.CloseTime.Time.Format("2006-01-02 15:04:05")
	}

	return item
}

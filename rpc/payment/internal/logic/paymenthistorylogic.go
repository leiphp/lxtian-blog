package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"lxtian-blog/common/pkg/model"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type PaymentHistoryLogic struct {
	*BaseLogic
}

func NewPaymentHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentHistoryLogic {
	return &PaymentHistoryLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
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

	offset := (in.Page - 1) * in.PageSize

	var paymentOrders []*model.PaymentOrder
	var err error

	// 根据查询条件获取支付记录
	if in.UserId > 0 {
		// 按用户ID查询
		paymentOrders, err = l.svcCtx.PaymentModel.FindPaymentOrdersByUserId(l.ctx, in.UserId, int(offset), int(in.PageSize))
		if err != nil {
			l.Errorf("Failed to find payment orders by user_id: %v", err)
			return nil, fmt.Errorf("failed to find payment orders by user_id: %w", err)
		}
		_, err = l.svcCtx.PaymentModel.CountPaymentOrdersByUserId(l.ctx, in.UserId)
		if err != nil {
			l.Errorf("Failed to count payment orders by user_id: %v", err)
		}
	} else if in.PaymentStatus != "" {
		// 按支付状态查询
		paymentOrders, err = l.svcCtx.PaymentModel.FindPaymentOrdersByStatus(l.ctx, in.PaymentStatus, int(offset), int(in.PageSize))
		if err != nil {
			l.Errorf("Failed to find payment orders by status: %v", err)
			return nil, fmt.Errorf("failed to find payment orders by status: %w", err)
		}
		_, err = l.svcCtx.PaymentModel.CountPaymentOrdersByStatus(l.ctx, in.PaymentStatus)
		if err != nil {
			l.Errorf("Failed to count payment orders by status: %v", err)
		}
	} else {
		// 查询所有记录
		paymentOrders, err = l.svcCtx.PaymentModel.FindPaymentOrdersByStatus(l.ctx, "", int(offset), int(in.PageSize))
		if err != nil {
			l.Errorf("Failed to find all payment orders: %v", err)
			return nil, fmt.Errorf("failed to find all payment orders: %w", err)
		}
		// 这里应该有一个查询所有记录总数的方法，暂时使用状态查询
		_, err = l.svcCtx.PaymentModel.CountPaymentOrdersByStatus(l.ctx, "")
		if err != nil {
			l.Errorf("Failed to count all payment orders: %v", err)
		}
	}

	// 时间过滤
	if in.StartTime != "" || in.EndTime != "" {
		paymentOrders = l.filterByTimeRange(paymentOrders, in.StartTime, in.EndTime)
	}

	// 订单ID过滤
	if in.OrderId != "" {
		paymentOrders = l.filterByOrderId(paymentOrders, in.OrderId)
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
		Total:    uint64(len(listData)), // 使用过滤后的数量
		List:     string(listJson),
	}, nil
}

// 按时间范围过滤
func (l *PaymentHistoryLogic) filterByTimeRange(orders []*model.PaymentOrder, startTime, endTime string) []*model.PaymentOrder {
	var filtered []*model.PaymentOrder

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

// 按订单ID过滤
func (l *PaymentHistoryLogic) filterByOrderId(orders []*model.PaymentOrder, orderId string) []*model.PaymentOrder {
	var filtered []*model.PaymentOrder

	for _, order := range orders {
		if order.OrderId == orderId {
			filtered = append(filtered, order)
		}
	}

	return filtered
}

// 构建支付订单项
func (l *PaymentHistoryLogic) buildPaymentOrderItem(order *model.PaymentOrder) map[string]interface{} {
	item := map[string]interface{}{
		"id":           order.ID,
		"payment_id":   order.PaymentId,
		"order_id":     order.OrderId,
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
		"client_ip":    order.ClientIP,
		"created_at":   order.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at":   order.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if order.BuyerUserId != "" {
		item["buyer_user_id"] = order.BuyerUserId
	}
	if order.BuyerLogonId != "" {
		item["buyer_logon_id"] = order.BuyerLogonId
	}
	if order.ReceiptAmount > 0 {
		item["receipt_amount"] = order.ReceiptAmount
	}
	if order.GmtPayment != nil {
		item["gmt_payment"] = order.GmtPayment.Format("2006-01-02 15:04:05")
	}
	if order.GmtClose != nil {
		item["gmt_close"] = order.GmtClose.Format("2006-01-02 15:04:05")
	}

	return item
}

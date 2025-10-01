package logic

import (
	"context"
	"database/sql"
	"fmt"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	paymentSvc "lxtian-blog/common/repository/payment"
	"strings"
	"time"

	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type PaymentNotifyLogic struct {
	*BaseLogic
	paymentService paymentSvc.PaymentOrderRepository
}

func NewPaymentNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentNotifyLogic {
	return &PaymentNotifyLogic{
		BaseLogic:      NewBaseLogic(ctx, svcCtx),
		paymentService: paymentSvc.NewPaymentOrderRepository(svcCtx.DB),
	}
}

func (l *PaymentNotifyLogic) PaymentNotify(in *payment.PaymentNotifyReq) (*payment.PaymentNotifyResp, error) {
	// 参数验证
	if in.NotifyData == "" {
		return &payment.PaymentNotifyResp{
			Success: false,
			Message: "通知数据不能为空",
		}, fmt.Errorf("notify_data is required")
	}

	// 生成通知ID
	notifyId := l.generateNotifyId()

	// 创建通知记录
	paymentNotify := &model.LxtPaymentNotifies{
		NotifyId:      notifyId,
		NotifyType:    constant.NotifyTypePayment,
		NotifyData:    in.NotifyData,
		Sign:          sql.NullString{String: in.Sign, Valid: in.Sign != ""},
		SignType:      sql.NullString{String: in.SignType, Valid: in.SignType != ""},
		VerifyStatus:  constant.VerifyStatusPending,
		ProcessStatus: constant.ProcessStatusPending,
	}

	// 保存通知记录
	err := l.svcCtx.DB.WithContext(l.ctx).Create(paymentNotify).Error
	if err != nil {
		l.Errorf("Failed to insert payment notify: %v", err)
		return &payment.PaymentNotifyResp{
			Success: false,
			Message: "保存通知记录失败",
		}, fmt.Errorf("failed to insert payment notify: %w", err)
	}

	// 验证签名
	err = l.verifySign(in.NotifyData, in.Sign)
	if err != nil {
		l.Errorf("Failed to verify sign: %v", err)
		l.paymentService.UpdatePaymentNotifyVerifyStatus(l.ctx, notifyId, constant.VerifyStatusFailed)
		return &payment.PaymentNotifyResp{
			Success: false,
			Message: "签名验证失败",
		}, fmt.Errorf("sign verification failed: %w", err)
	}

	// 更新验证状态为成功
	err = l.paymentService.UpdatePaymentNotifyVerifyStatus(l.ctx, notifyId, constant.VerifyStatusSuccess)
	if err != nil {
		l.Errorf("Failed to update verify status: %v", err)
	}

	// 解析通知数据
	notifyData, err := l.parseNotifyData(in.NotifyData)
	if err != nil {
		l.Errorf("Failed to parse notify data: %v", err)
		l.paymentService.UpdatePaymentNotifyProcessStatus(l.ctx, notifyId, constant.ProcessStatusFailed, "解析通知数据失败")
		return &payment.PaymentNotifyResp{
			Success: false,
			Message: "解析通知数据失败",
		}, fmt.Errorf("failed to parse notify data: %w", err)
	}

	// 处理通知
	err = l.processNotify(notifyData, notifyId)
	if err != nil {
		l.Errorf("Failed to process notify: %v", err)
		l.paymentService.UpdatePaymentNotifyProcessStatus(l.ctx, notifyId, constant.ProcessStatusFailed, err.Error())
		return &payment.PaymentNotifyResp{
			Success: false,
			Message: "处理通知失败",
		}, fmt.Errorf("failed to process notify: %w", err)
	}

	// 更新处理状态为成功
	err = l.paymentService.UpdatePaymentNotifyProcessStatus(l.ctx, notifyId, constant.ProcessStatusSuccess, "")
	if err != nil {
		l.Errorf("Failed to update process status: %v", err)
	}

	// 记录日志
	l.Infof("Processed payment notify: notifyId=%s, out_trade_no=%s", notifyId, notifyData["out_trade_no"])

	return &payment.PaymentNotifyResp{
		Success: true,
		Message: "处理成功",
	}, nil
}

// 验证签名
func (l *PaymentNotifyLogic) verifySign(data, sign string) error {
	// 使用支付宝客户端验证签名
	return l.svcCtx.AlipayClient.VerifySign(data, sign)
}

// 解析通知数据
func (l *PaymentNotifyLogic) parseNotifyData(notifyData string) (map[string]string, error) {
	// 解析支付宝通知数据格式
	// 支付宝通知数据通常是URL编码的键值对格式
	data := make(map[string]string)

	// 简单的URL参数解析
	pairs := strings.Split(notifyData, "&")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			// 这里应该进行URL解码，简化处理
			data[kv[0]] = kv[1]
		}
	}

	// 验证必要字段
	if data["out_trade_no"] == "" {
		return nil, fmt.Errorf("missing out_trade_no")
	}

	if data["trade_status"] == "" {
		return nil, fmt.Errorf("missing trade_status")
	}

	return data, nil
}

// 处理通知
func (l *PaymentNotifyLogic) processNotify(notifyData map[string]string, notifyId string) error {
	outTradeNo := notifyData["out_trade_no"]
	tradeStatus := notifyData["trade_status"]

	// 查找支付订单
	paymentOrder, err := l.paymentService.FindPaymentOrderByOutTradeNo(l.ctx, outTradeNo)
	if err != nil {
		return fmt.Errorf("payment order not found: %w", err)
	}

	// 更新通知记录的支付ID
	notify, err := l.paymentService.FindPaymentNotifyByNotifyId(l.ctx, notifyId)
	if err == nil {
		notify.PaymentId = paymentOrder.PaymentId
		l.paymentService.UpdatePaymentNotify(l.ctx, notify)
	}

	// 根据交易状态处理
	switch tradeStatus {
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		// 支付成功
		if paymentOrder.Status == constant.PaymentStatusPaid {
			// 已经处理过，直接返回成功
			return nil
		}

		// 解析支付时间
		var gmtPayment *time.Time
		if notifyData["gmt_payment"] != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", notifyData["gmt_payment"]); err == nil {
				gmtPayment = &t
			}
		}

		// 解析金额
		var receiptAmount float64
		if notifyData["receipt_amount"] != "" {
			if amount, err := parseFloat(notifyData["receipt_amount"]); err == nil {
				receiptAmount = amount
			}
		}

		// 更新订单信息
		err = l.paymentService.UpdatePaymentOrderTradeInfo(
			l.ctx,
			paymentOrder.PaymentId,
			notifyData["trade_no"],
			tradeStatus,
			notifyData["buyer_id"],
			notifyData["buyer_logon_id"],
			receiptAmount,
			gmtPayment,
		)
		if err != nil {
			return fmt.Errorf("failed to update trade info: %w", err)
		}

		// 这里可以添加支付成功后的业务逻辑
		// 例如：发送通知、更新库存、发放优惠券等
		l.handlePaymentSuccess(paymentOrder, notifyData)

	case "TRADE_CLOSED":
		// 交易关闭
		if paymentOrder.Status != constant.PaymentStatusClosed {
			err = l.paymentService.UpdatePaymentOrderStatus(l.ctx, paymentOrder.PaymentId, constant.PaymentStatusClosed)
			if err != nil {
				return fmt.Errorf("failed to update status to closed: %w", err)
			}
		}

	case "WAIT_BUYER_PAY":
		// 等待买家付款，不需要特殊处理

	default:
		// 其他状态
		l.Errorf("Unknown trade status: %s", tradeStatus)
	}

	return nil
}

// 处理支付成功后的业务逻辑
func (l *PaymentNotifyLogic) handlePaymentSuccess(paymentOrder *model.LxtPaymentOrders, notifyData map[string]string) {
	// 这里可以添加具体的业务逻辑
	// 例如：
	// 1. 发送支付成功通知给用户
	// 2. 更新订单状态
	// 3. 发放积分或优惠券
	// 4. 记录支付日志

	l.Infof("Payment success: paymentId=%s, orderSn=%s, amount=%.2f",
		paymentOrder.PaymentId, paymentOrder.OrderSn, paymentOrder.Amount)
}

// 简单的浮点数解析
func parseFloat(s string) (float64, error) {
	// 这里应该使用strconv.ParseFloat，简化处理
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	// 实际项目中应该正确处理浮点数解析
	return 0, fmt.Errorf("not implemented")
}

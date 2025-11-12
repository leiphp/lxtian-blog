package logic

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	paymentSvc "lxtian-blog/common/repository/payment_repo"

	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"

	"gorm.io/gorm"
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
	if strings.TrimSpace(in.NotifyData) == "" {
		return &payment.PaymentNotifyResp{
			Success: false,
			Message: "通知数据不能为空",
		}, fmt.Errorf("notify_data is required")
	}

	// 生成通知ID
	notifyId := l.generateNotifyId()

	// 创建通知记录
	paymentNotify := &model.LxtPaymentNotify{
		NotifyID:      notifyId,
		NotifyType:    constant.NotifyTypePayment,
		NotifyData:    in.NotifyData,
		Sign:          &in.Sign,
		SignType:      &in.SignType,
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
	if strings.TrimSpace(sign) == "" {
		return fmt.Errorf("sign is empty")
	}

	// 构建待验签内容
	signContent, err := buildSignContent(data)
	if err != nil {
		return fmt.Errorf("failed to build sign content: %w", err)
	}
	if signContent == "" {
		return fmt.Errorf("sign content is empty")
	}

	l.Infof("Verifying sign, content length: %d", len(signContent))
	l.Infof("Verifying sign, sign length: %d", len(sign))
	l.Infof("Verifying sign, sign (first 50 chars): %s", safeSubstring(sign, 50))

	// 使用支付宝客户端验证签名
	err = l.svcCtx.AlipayClient.VerifySign(signContent, sign)
	if err != nil {
		l.Errorf("Alipay sign verification failed: %v", err)
		return fmt.Errorf("alipay sign verification failed: %w", err)
	}

	l.Infof("Alipay sign verification success")
	return nil
}

// 安全截取字符串
func safeSubstring(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// buildSignContent 构建待验签字符串
func buildSignContent(rawData string) (string, error) {
	if strings.TrimSpace(rawData) == "" {
		return "", fmt.Errorf("raw notify data is empty")
	}

	segments := strings.Split(rawData, "&")
	if len(segments) == 0 {
		return "", fmt.Errorf("invalid notify data")
	}

	values := make(map[string][]string)
	for _, segment := range segments {
		if segment == "" {
			continue
		}

		kv := strings.SplitN(segment, "=", 2)
		key := kv[0]
		var value string
		if len(kv) == 2 {
			value = kv[1]
		}

		if key == "sign" || key == "sign_type" {
			continue
		}

		decodedKey, err := url.QueryUnescape(key)
		if err != nil {
			decodedKey = key
		}

		decodedValue, err := url.QueryUnescape(value)
		if err != nil {
			decodedValue = value
		}

		values[decodedKey] = append(values[decodedKey], decodedValue)
	}

	if len(values) == 0 {
		return "", fmt.Errorf("no parameters available for sign")
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	first := true
	for _, key := range keys {
		for _, value := range values[key] {
			if !first {
				builder.WriteByte('&')
			} else {
				first = false
			}
			builder.WriteString(key)
			builder.WriteByte('=')
			builder.WriteString(value)
		}
	}

	return builder.String(), nil
}

// 解析通知数据
func (l *PaymentNotifyLogic) parseNotifyData(notifyData string) (map[string]string, error) {
	data := make(map[string]string)

	pairs := strings.Split(notifyData, "&")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key, err := url.QueryUnescape(kv[0])
		if err != nil {
			key = kv[0]
		}
		value, err := url.QueryUnescape(kv[1])
		if err != nil {
			value = kv[1]
		}

		data[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	// 验证必要字段
	if data["out_trade_no"] == "" {
		return nil, fmt.Errorf("missing out_trade_no")
	}

	if data["trade_status"] == "" {
		return nil, fmt.Errorf("missing trade_status")
	}

	// 统一交易状态格式
	data["trade_status"] = strings.ToUpper(data["trade_status"])

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
		notify.PaymentID = paymentOrder.PaymentID
		if err := l.paymentService.UpdatePaymentNotify(l.ctx, notify); err != nil {
			l.Errorf("Failed to update payment notify record: %v", err)
		}
	}

	// 根据交易状态处理
	switch tradeStatus {
	case constant.TradeStatusSuccess, constant.TradeStatusFinished:
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
			} else {
				l.Errorf("Failed to parse gmt_payment '%s': %v", notifyData["gmt_payment"], err)
			}
		}

		// 解析金额
		var receiptAmount float64
		if notifyData["receipt_amount"] != "" {
			if amount, err := parseFloat(notifyData["receipt_amount"]); err == nil {
				receiptAmount = amount
			} else {
				l.Errorf("Failed to parse receipt_amount '%s': %v", notifyData["receipt_amount"], err)
			}
		}

		// 更新订单信息
		err = l.paymentService.UpdatePaymentOrderTradeInfo(
			l.ctx,
			paymentOrder.PaymentID,
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

		// 显式更新订单状态为已支付，确保状态同步
		//if err := l.paymentService.UpdatePaymentOrderStatus(l.ctx, paymentOrder.PaymentID, constant.PaymentStatusPaid); err != nil {
		//	return fmt.Errorf("failed to update payment order status: %w", err)
		//}

		// 同步更新内存对象，便于后续逻辑使用
		paymentOrder.TradeNo = notifyData["trade_no"]
		paymentOrder.TradeStatus = tradeStatus
		paymentOrder.BuyerUserID = notifyData["buyer_id"]
		paymentOrder.BuyerLogonID = notifyData["buyer_logon_id"]
		paymentOrder.ReceiptAmount = receiptAmount
		paymentOrder.Status = constant.PaymentStatusPaid
		if gmtPayment != nil {
			paymentOrder.PayTime = gmtPayment
		}

		// 支付成功后的业务逻辑
		if err := l.handlePaymentSuccess(paymentOrder, notifyData); err != nil {
			return err
		}

	case constant.TradeStatusClosed:
		// 交易关闭
		if paymentOrder.Status != constant.PaymentStatusClosed {
			err = l.paymentService.UpdatePaymentOrderStatus(l.ctx, paymentOrder.PaymentID, constant.PaymentStatusClosed)
			if err != nil {
				return fmt.Errorf("failed to update status to closed: %w", err)
			}
		}

	case constant.TradeStatusWaitBuyerPay:
		// 等待买家付款，不需要特殊处理

	default:
		// 其他状态
		l.Errorf("Unknown trade status: %s", tradeStatus)
	}

	return nil
}

// 处理支付成功后的业务逻辑
func (l *PaymentNotifyLogic) handlePaymentSuccess(paymentOrder *model.LxtPaymentOrder, notifyData map[string]string) error {
	l.Infof("Payment success: paymentId=%s, orderSn=%s, amount=%.2f",
		paymentOrder.PaymentID, paymentOrder.OrderSn, paymentOrder.Amount)

	// 如果是会员购买订单，自动开通或续费会员
	if paymentOrder.BuyType == constant.BuyTypeMembership {
		if err := l.activateMembership(paymentOrder); err != nil {
			l.Errorf("Failed to activate membership for paymentId=%s: %v", paymentOrder.PaymentID, err)
			return fmt.Errorf("activate membership failed: %w", err)
		}
	}

	return nil
}

// activateMembership 为会员订单开通或续费会员
func (l *PaymentNotifyLogic) activateMembership(paymentOrder *model.LxtPaymentOrder) error {
	if paymentOrder == nil {
		return errors.New("payment order is nil")
	}

	return l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		var membershipType model.LxtUserMembershipType
		if err := tx.Where("id = ?", paymentOrder.GoodsID).First(&membershipType).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("membership type not found for goods_id=%d", paymentOrder.GoodsID)
			}
			return fmt.Errorf("query membership type failed: %w", err)
		}

		daysToAdd := int(membershipType.Days)
		if daysToAdd <= 0 {
			return fmt.Errorf("invalid membership days for type %d", membershipType.ID)
		}

		now := time.Now()

		var membership model.LxtUserMembership
		err := tx.Where("user_id = ?", paymentOrder.UserID).First(&membership).Error

		var (
			beforeStart      *time.Time
			beforeEnd        *time.Time
			fromTypeID       *int64
			remainingDaysPtr *int32
			renewalType      int32 = constant.MembershipRenewalTypeRenewal
		)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			membership = model.LxtUserMembership{
				UserID:           paymentOrder.UserID,
				MembershipTypeID: membershipType.ID,
				StartTime:        now,
				EndTime:          now.AddDate(0, 0, daysToAdd),
				IsActive:         1,
				TotalDays:        int32(daysToAdd),
				Level:            calculateMembershipLevel(int(daysToAdd)),
			}

			if err := tx.Create(&membership).Error; err != nil {
				return fmt.Errorf("create membership failed: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("query membership failed: %w", err)
		} else {
			origStart := membership.StartTime
			origEnd := membership.EndTime
			beforeStart = &origStart
			beforeEnd = &origEnd
			fromType := membership.MembershipTypeID
			fromTypeID = &fromType

			if membership.MembershipTypeID != membershipType.ID {
				renewalType = constant.MembershipRenewalTypeUpgrade
			}

			if membership.EndTime.Before(now) {
				membership.StartTime = now
				membership.EndTime = now.AddDate(0, 0, daysToAdd)
			} else {
				membership.EndTime = membership.EndTime.AddDate(0, 0, daysToAdd)
				remaining := int32(math.Ceil(origEnd.Sub(now).Hours() / 24))
				if remaining < 0 {
					remaining = 0
				}
				remainingCopy := remaining
				remainingDaysPtr = &remainingCopy
			}

			membership.MembershipTypeID = membershipType.ID
			membership.IsActive = 1
			newTotal := membership.TotalDays + int32(daysToAdd)
			membership.TotalDays = newTotal
			membership.Level = calculateMembershipLevel(int(newTotal))

			if err := tx.Save(&membership).Error; err != nil {
				return fmt.Errorf("update membership failed: %w", err)
			}
		}

		orderID := paymentOrder.ID
		renewalRecord := &model.LxtUserMembershipRenewal{
			UserID:               paymentOrder.UserID,
			OrderID:              &orderID,
			FromMembershipTypeID: fromTypeID,
			ToMembershipTypeID:   membershipType.ID,
			RenewalType:          renewalType,
			BeforeStartTime:      beforeStart,
			BeforeEndTime:        beforeEnd,
			AfterStartTime:       membership.StartTime,
			AfterEndTime:         membership.EndTime,
			RemainingDays:        remainingDaysPtr,
			CalculatedDays:       membership.TotalDays,
			Amount:               paymentOrder.Amount,
		}

		if err := tx.Create(renewalRecord).Error; err != nil {
			return fmt.Errorf("create membership renewal failed: %w", err)
		}

		return nil
	})
}

// calculateMembershipLevel 按累计天数计算会员等级
func calculateMembershipLevel(totalDays int) int32 {
	switch {
	case totalDays <= 90:
		return 1
	case totalDays <= 180:
		return 2
	case totalDays <= 365:
		return 3
	case totalDays <= 730:
		return 4
	default:
		return 5
	}
}

// 简单的浮点数解析
func parseFloat(s string) (float64, error) {
	if strings.TrimSpace(s) == "" {
		return 0, fmt.Errorf("empty string")
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

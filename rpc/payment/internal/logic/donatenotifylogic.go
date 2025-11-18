package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	redisutil "lxtian-blog/common/pkg/redis"
	paymentSvc "lxtian-blog/common/repository/payment_repo"
	"lxtian-blog/common/repository/web_repo"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
	"net/url"
	"strings"
	"time"
)

type DonateNotifyLogic struct {
	*BaseLogic
	paymentService paymentSvc.PaymentOrderRepository
	webService     web_repo.TxyOrderRepository
}

func NewDonateNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DonateNotifyLogic {
	return &DonateNotifyLogic{
		BaseLogic:      NewBaseLogic(ctx, svcCtx),
		paymentService: paymentSvc.NewPaymentOrderRepository(svcCtx.DB),
		webService:     web_repo.NewTxyOrderRepository(svcCtx.DB),
	}
}

// 捐赠支付回调通知处理
func (l *DonateNotifyLogic) DonateNotify(in *payment.DonateNotifyReq) (*payment.DonateNotifyResp, error) {
	// 参数验证
	if strings.TrimSpace(in.NotifyData) == "" {
		return &payment.DonateNotifyResp{
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
		return &payment.DonateNotifyResp{
			Success: false,
			Message: "保存通知记录失败",
		}, fmt.Errorf("failed to insert payment notify: %w", err)
	}

	// 验证签名
	err = l.verifySign(in.NotifyData, in.Sign)
	if err != nil {
		l.Errorf("Failed to verify sign: %v", err)
		l.paymentService.UpdatePaymentNotifyVerifyStatus(l.ctx, notifyId, constant.VerifyStatusFailed)
		return &payment.DonateNotifyResp{
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
		return &payment.DonateNotifyResp{
			Success: false,
			Message: "解析通知数据失败",
		}, fmt.Errorf("failed to parse notify data: %w", err)
	}

	// 处理通知
	err = l.processNotify(notifyData, notifyId)
	if err != nil {
		l.Errorf("Failed to process notify: %v", err)
		l.paymentService.UpdatePaymentNotifyProcessStatus(l.ctx, notifyId, constant.ProcessStatusFailed, err.Error())
		return &payment.DonateNotifyResp{
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

	return &payment.DonateNotifyResp{
		Success: true,
		Message: "处理成功",
	}, nil
}

// 验证签名
func (l *DonateNotifyLogic) verifySign(data, sign string) error {
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

	// 使用支付宝客户端验证签名
	err = l.svcCtx.AlipayClient.VerifySign(signContent, sign)
	if err != nil {
		l.Errorf("Alipay sign verification failed: %v", err)
		return fmt.Errorf("alipay sign verification failed: %w", err)
	}

	l.Infof("Alipay sign verification success")
	return nil
}

// 解析通知数据
func (l *DonateNotifyLogic) parseNotifyData(notifyData string) (map[string]string, error) {
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
func (l *DonateNotifyLogic) processNotify(notifyData map[string]string, notifyId string) error {
	outTradeNo := notifyData["out_trade_no"]
	tradeStatus := notifyData["trade_status"]

	// 从Redis获取支付订单
	txyOrder, err := l.getTxyOrderFromRedis(outTradeNo)
	if err != nil {
		return fmt.Errorf("payment order not found: %w", err)
	}

	// 更新通知记录的支付ID
	notify, err := l.paymentService.FindPaymentNotifyByNotifyId(l.ctx, notifyId)
	if err == nil {
		notify.PaymentID = txyOrder.PaymentID
		if err := l.paymentService.UpdatePaymentNotify(l.ctx, notify); err != nil {
			l.Errorf("Failed to update payment notify record: %v", err)
		}
	}

	// 根据交易状态处理
	switch tradeStatus {
	case constant.TradeStatusSuccess, constant.TradeStatusFinished:
		// 支付成功，插入订单状态为已支付
		txyOrder.Status = constant.PaymentStatusPaid
		txyOrder.UpdatedAt = time.Now()

		// 使用GORM保存支付订单到数据库
		err = l.svcCtx.DB.WithContext(l.ctx).Save(txyOrder).Error
		if err != nil {
			l.Errorf("Failed to update web_repo order: %v", err)
			return fmt.Errorf("更新捐赠支付订单失败: %w", err)
		}
		return nil

	case constant.TradeStatusClosed:
		// 交易关闭
		if txyOrder.Status != constant.PaymentStatusClosed {
			err = l.webService.UpdateStatus(l.ctx, txyOrder.PaymentID, constant.PaymentStatusClosed)
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

// getTxyOrderFromRedis 从Redis获取支付订单
func (l *DonateNotifyLogic) getTxyOrderFromRedis(outTradeNo string) (*model.TxyOrder, error) {
	if l.svcCtx.Rds == nil {
		return nil, fmt.Errorf("redis client is nil")
	}

	// 生成Redis key
	redisKey := redisutil.ReturnRedisKey(redisutil.DonatePendingOrderString, outTradeNo)

	// 从Redis获取订单数据
	orderJson, err := l.svcCtx.Rds.Get(redisKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get order from redis: %w", err)
	}

	if orderJson == "" {
		return nil, fmt.Errorf("order not found in redis")
	}

	// 反序列化为TxyOrder
	var txyOrder model.TxyOrder
	if err := json.Unmarshal([]byte(orderJson), &txyOrder); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order from redis: %w", err)
	}

	return &txyOrder, nil
}

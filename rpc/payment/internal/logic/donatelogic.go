package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
	redisutil "lxtian-blog/common/pkg/redis"
	"lxtian-blog/common/pkg/utils"
	"strconv"
	"time"

	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
)

type DonateLogic struct {
	*BaseLogic
}

func NewDonateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DonateLogic {
	return &DonateLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
	}
}

// 创建捐赠订单
func (l *DonateLogic) Donate(in *payment.DonateReq) (*payment.DonateResp, error) {
	// 参数验证
	if in.Amount <= 0 {
		return nil, fmt.Errorf("支付金额必须大于0")
	}

	if in.Subject == "" {
		return nil, fmt.Errorf("订单标题不能为空")
	}

	// 检查用户是否有待支付的捐赠订单
	// 无论是否登录，都要限制待支付订单数量为3单
	var userIdentifier string
	if in.UserId > 0 {
		// 已登录用户，使用UserId标识
		userIdentifier = fmt.Sprintf("user:%d", in.UserId)
	} else {
		// 未登录用户，使用IP地址标识
		if in.ClientIp == "" {
			return nil, fmt.Errorf("未登录用户必须提供客户端IP")
		}
		userIdentifier = fmt.Sprintf("ip:%s", in.ClientIp)
	}

	// 获取当前待支付订单数量
	count, err := l.getDonatePendingCount(userIdentifier)
	if err != nil {
		l.Errorf("Failed to get pending donate orders count: %v", err)
		return nil, fmt.Errorf("检查待支付订单失败: %w", err)
	}

	// 如果已有3单未支付订单，不允许再创建
	if count >= 3 {
		return nil, fmt.Errorf("您已有3单待支付的捐赠订单，请勿频繁创建新订单")
	}

	// 生成订单ID、支付ID和商户订单号
	paymentId := l.generatePaymentId()
	outTradeNo := fmt.Sprintf("NO%d", utils.Snowflake())
	orderSn := fmt.Sprintf("%d", utils.Snowflake())

	// 产品码：电脑网站支付固定使用 FAST_INSTANT_TRADE_PAY
	if in.ProductCode == "" {
		in.ProductCode = "FAST_INSTANT_TRADE_PAY"
	}

	if in.Timeout == "" {
		in.Timeout = "30m"
	}
	// 设置默认值（需要重新生成protobuf后启用）
	if in.PayType == 0 {
		in.PayType = 1 // 支付宝
	}
	// 设置默认值（需要重新生成protobuf后启用）
	if in.BuyType == 0 {
		in.BuyType = 1 // 默认捐赠
	}

	// 2. 创建支付订单记录（txy_orders表）
	paymentOrder := &model.TxyOrder{
		PaymentID:  paymentId,
		OrderSn:    orderSn,
		OutTradeNo: outTradeNo,
		UserID:     int32(in.UserId),
		Amount:     in.Amount,
		Subject:    in.Subject,
		Status:     constant.PaymentStatusPending,
		PayType:    int32(in.PayType),
		Remark:     in.Remark,
		IP:         in.ClientIp,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 使用GORM保存支付订单到数据库 todo 待支付不入库
	//err = l.svcCtx.DB.WithContext(l.ctx).Create(paymentOrder).Error
	//if err != nil {
	//	l.Errorf("Failed to insert payment_repo order: %v", err)
	//	return nil, fmt.Errorf("创建支付订单失败: %w", err)
	//}

	// 将paymentOrder数据存储到Redis，key=donate:order:{userIdentifier}:{outTradeNo}
	err = l.savePaymentOrderToRedis(paymentOrder, userIdentifier, outTradeNo)
	if err != nil {
		l.Errorf("Failed to save payment order to redis: %v", err)
		// 这里不返回错误，因为订单已经创建成功，只是Redis记录失败
	}

	// 3. 调用支付宝API创建支付订单
	// 金额转换为字符串，保留2位小数
	amountStr := strconv.FormatFloat(in.Amount, 'f', 2, 64)

	// 超时时间格式：30m, 1h, 1d 等，默认30m
	timeout := in.Timeout
	if timeout == "" || !isValidTimeout(timeout) {
		timeout = "30m"
	}

	alipayReq := &alipay.TradeCreateRequest{
		OutTradeNo:  outTradeNo,
		TotalAmount: amountStr,
		Subject:     in.Subject,
		Body:        in.Body,
		ProductCode: in.ProductCode,
		Timeout:     timeout,
		ReturnUrl:   in.ReturnUrl,
	}

	// 调用支付宝API创建支付URL
	payUrl, err := l.svcCtx.AlipayClient.CreatePayment(alipayReq)
	if err != nil {
		l.Errorf("Failed to create alipay payment_repo: %v", err)
		// 使用GORM更新订单状态为失败
		l.svcCtx.DB.WithContext(l.ctx).Model(&model.LxtPaymentOrder{}).
			Where("payment_id = ?", paymentId).
			Update("status", constant.VerifyStatusFailed)
		return nil, fmt.Errorf("创建支付订单失败: %w", err)
	}

	// 记录日志
	l.Infof("Created payment_repo order: paymentId=%s, orderSn=%s, outTradeNo=%s, amount=%.2f, payUrl=%s",
		paymentId, orderSn, outTradeNo, in.Amount, payUrl)

	return &payment.DonateResp{
		PaymentId:  paymentId,
		OutTradeNo: outTradeNo,
		OrderSn:    orderSn,
		PayUrl:     payUrl,
	}, nil
}

// getDonatePendingCount 获取待支付捐赠订单数量
// userIdentifier格式：user:{userId} 或 ip:{ip}
// 通过检查Set中的订单数量来统计（Set的key为blog:donate:pending:{userIdentifier}）
// 当订单过期时，需要从Set中移除对应的订单号
func (l *DonateLogic) getDonatePendingCount(userIdentifier string) (int, error) {
	if l.svcCtx.Rds == nil {
		return 0, nil
	}

	// 使用Set来存储订单列表，key格式：blog:donate:pending:{userIdentifier}
	// 必须与savePaymentOrderToRedis中使用的key格式一致
	setKey := redisutil.ReturnRedisKey(redisutil.DonatePendingOrderSet, userIdentifier)

	// 获取Set的大小（订单数量）
	count, err := l.svcCtx.Rds.ScardCtx(l.ctx, setKey)
	if err != nil {
		// 如果获取失败，返回0
		l.Errorf("Failed to get order count from set: %v", err)
		return 0, nil
	}

	// 清理Set中已过期的订单（通过检查订单key是否存在）
	// 这样可以确保计数准确
	l.cleanExpiredOrdersFromSet(setKey)

	// 重新获取清理后的数量
	count, err = l.svcCtx.Rds.ScardCtx(l.ctx, setKey)
	if err != nil {
		return 0, nil
	}

	return int(count), nil
}

// cleanExpiredOrdersFromSet 清理Set中已过期的订单
func (l *DonateLogic) cleanExpiredOrdersFromSet(setKey string) {
	// 获取Set中的所有订单号
	orderNos, err := l.svcCtx.Rds.SmembersCtx(l.ctx, setKey)
	if err != nil {
		l.Errorf("Failed to get order nos from set: %v", err)
		return
	}

	// 检查每个订单是否存在
	for _, orderNo := range orderNos {
		orderKey := redisutil.ReturnRedisKey(redisutil.DonatePendingOrderString, orderNo)
		exists, err := l.svcCtx.Rds.ExistsCtx(l.ctx, orderKey)
		if err != nil {
			continue
		}
		// 如果订单不存在（已过期），从Set中移除
		if !exists {
			_, _ = l.svcCtx.Rds.SremCtx(l.ctx, setKey, orderNo)
		}
	}
}

// savePaymentOrderToRedis 将paymentOrder数据存储到Redis
// key格式：donate:order:{outTradeNo}
// 同时将订单号添加到Set中，用于统计订单数量
// 订单过期后会自动从Redis删除，但需要从Set中手动移除（通过cleanExpiredOrdersFromSet）
func (l *DonateLogic) savePaymentOrderToRedis(paymentOrder *model.TxyOrder, userIdentifier, outTradeNo string) error {
	if l.svcCtx.Rds == nil {
		return nil
	}

	// 将paymentOrder序列化为JSON
	orderJsonBytes, err := json.Marshal(paymentOrder)
	if err != nil {
		return fmt.Errorf("序列化订单数据失败: %w", err)
	}

	// key格式：donate:order:{outTradeNo}
	redisKey := redisutil.ReturnRedisKey(redisutil.DonatePendingOrderString, outTradeNo)

	// 保存订单数据到Redis，设置30分钟过期时间（1800秒）
	// 过期后Redis会自动删除
	err = l.svcCtx.Rds.SetexCtx(l.ctx, redisKey, string(orderJsonBytes), 1800)
	if err != nil {
		return fmt.Errorf("保存订单到Redis失败: %w", err)
	}

	// 将订单号添加到Set中，用于统计订单数量
	// Set的key格式：donate:pending:{userIdentifier}
	setKey := redisutil.ReturnRedisKey(redisutil.DonatePendingOrderSet, userIdentifier)
	_, err = l.svcCtx.Rds.SaddCtx(l.ctx, setKey, outTradeNo)
	if err != nil {
		// Set操作失败，记录日志但不返回错误
		l.Errorf("Failed to add order to set: %v", err)
	} else {
		// 设置Set的过期时间为30分钟（1800秒）
		// 注意：Set的过期时间应该比订单的过期时间稍长，以确保清理逻辑能够执行
		err = l.svcCtx.Rds.ExpireCtx(l.ctx, setKey, 1800)
		if err != nil {
			l.Errorf("Failed to set set key expire: %v", err)
		}
	}

	return nil
}

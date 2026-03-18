package logic

import (
	"context"
	"fmt"
	"github.com/leiphp/unit-go-sdk/pkg/gconv"
	"lxtian-blog/common/constant"
	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/alipay"
	redisutil "lxtian-blog/common/pkg/redis"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/common/repository/user_repo"
	"lxtian-blog/rpc/payment/internal/svc"
	"lxtian-blog/rpc/payment/pb/payment"
	"strconv"
	"time"
)

type CreatePaymentLogic struct {
	*BaseLogic
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		BaseLogic: NewBaseLogic(ctx, svcCtx),
	}
}

func (l *CreatePaymentLogic) CreatePayment(in *payment.CreatePaymentReq) (*payment.CreatePaymentResp, error) {
	var err error
	// 参数验证
	if in.Amount <= 0 {
		return nil, fmt.Errorf("支付金额必须大于0")
	}

	if in.Subject == "" {
		return nil, fmt.Errorf("订单标题不能为空")
	}

	if in.UserId == 0 {
		return nil, fmt.Errorf("用户ID不能为空")
	}

	// 使用 Redis 控制每日「待支付订单」上限：
	// 每个用户每天最多允许 3 个“仍处于待支付状态”的订单。
	// 这里只是预检查当前计数，真正的 +1 放在订单创建成功之后。
	if l.svcCtx.Rds != nil {
		today := time.Now()
		dateStr := today.Format("2006-01-02")
		dailyKey := fmt.Sprintf("%spayment:pending:count:%d:%s", redisutil.KeyPrefix, in.UserId, dateStr)

		var current int64
		if val, err := l.svcCtx.Rds.GetCtx(l.ctx, dailyKey); err == nil && val != "" {
			if parsed, perr := strconv.ParseInt(val, 10, 64); perr == nil {
				current = parsed
			}
		}
		if current >= 3 {
			return nil, fmt.Errorf("今日待支付订单数量已达上限（3个），请处理现有订单后再创建")
		}
	}

	// 检查用户是否购买会员
	if in.VipId > 0 {
		// 查询并处理会员信息（优先从 Redis 获取，未命中再查 DB）
		membershipRepo := user_repo.NewUserMembershipRepository(l.svcCtx.DB, l.svcCtx.Rds)
		membershipInfo, err := membershipRepo.GetActiveMembershipByUserId(l.ctx, int64(in.UserId))
		if err != nil {
			l.Errorf("获取用户会员信息失败: %v", err)
			// 会员信息查询失败不影响用户信息返回，继续执行
			membershipInfo = nil
		}
		if membershipInfo != nil && gconv.Bool(membershipInfo["is_valid"]) {
			if in.VipId < gconv.Int64(membershipInfo["type_id"]) {
				return nil, fmt.Errorf("会员只允许升级和续费，不能降级")
			}
		}
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
		in.PayType = 1 // 默认直接消费
	}
	// 设置默认值（需要重新生成protobuf后启用）
	if in.BuyType == 0 {
		in.PayType = 3 // 默认直接消费
	}
	if in.Quantity == 0 {
		in.Quantity = 1 // 默认数量为1
	}

	// 2. 创建支付订单记录（lxt_payment_orders表）
	paymentOrder := &model.LxtPaymentOrder{
		GoodsID:     int32(in.GoodsId),
		Quantity:    int32(in.Quantity),
		VipID:       int32(in.VipId),
		PaymentID:   paymentId,
		OrderSn:     orderSn,
		OutTradeNo:  outTradeNo,
		UserID:      int64(in.UserId),
		Amount:      in.Amount,
		Subject:     in.Subject,
		Status:      constant.PaymentStatusPending,
		ProductCode: &in.ProductCode,
		PayType:     int32(in.PayType),
		BuyType:     int32(in.BuyType),
		// 处理可空字段
		Body:      &in.Body,
		ReturnURL: in.ReturnUrl,
		NotifyURL: in.NotifyUrl,
		Timeout:   in.Timeout,
		ClientIP:  in.ClientIp,
		Remark:    in.Remark,
	}

	// 使用GORM保存支付订单到数据库
	err = l.svcCtx.DB.WithContext(l.ctx).Create(paymentOrder).Error
	if err != nil {
		l.Errorf("Failed to insert payment_repo order: %v", err)
		return nil, fmt.Errorf("创建支付订单失败: %w", err)
	}

	// 订单已成功写入数据库，此时将「待支付订单计数」+1
	if l.svcCtx.Rds != nil {
		today := paymentOrder.CreatedAt
		dateStr := today.Format("2006-01-02")
		dailyKey := fmt.Sprintf("%spayment:pending:count:%d:%s", redisutil.KeyPrefix, paymentOrder.UserID, dateStr)

		dailyCount, rerr := l.svcCtx.Rds.IncrCtx(l.ctx, dailyKey)
		if rerr != nil {
			l.Errorf("failed to incr daily pending order count after create, key=%s, err=%v", dailyKey, rerr)
		} else if dailyCount == 1 {
			// 首次创建时，设置 key 在当天 23:59:59 过期
			endOfDay := time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 0, today.Location())
			ttl := int(endOfDay.Sub(time.Now()).Seconds())
			if ttl > 0 {
				if err := l.svcCtx.Rds.ExpireCtx(l.ctx, dailyKey, ttl); err != nil {
					l.Errorf("failed to set expire for key=%s, err=%v", dailyKey, err)
				}
			}
		}
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

	return &payment.CreatePaymentResp{
		PaymentId:  paymentId,
		OutTradeNo: outTradeNo,
		OrderSn:    orderSn,
		PayUrl:     payUrl,
	}, nil
}

// 验证超时时间格式是否有效
func isValidTimeout(timeout string) bool {
	// 支持的格式：30m, 1h, 1d, 1c（c表示天）
	if len(timeout) < 2 {
		return false
	}

	// 检查是否以数字开头，以 m/h/d/c 结尾
	lastChar := timeout[len(timeout)-1]
	return (lastChar == 'm' || lastChar == 'h' || lastChar == 'd' || lastChar == 'c') &&
		timeout[:len(timeout)-1] != "" &&
		isNumeric(timeout[:len(timeout)-1])
}

// 检查字符串是否为数字
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

package logic

import (
	"context"
	"time"

	"lxtian-blog/rpc/payment/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BaseLogic {
	return &BaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成支付ID
func (l *BaseLogic) generatePaymentId() string {
	return "PAY_" + time.Now().Format("20060102150405") + "_" + l.generateRandomString(8)
}

// 生成退款ID
func (l *BaseLogic) generateRefundId() string {
	return "REF_" + time.Now().Format("20060102150405") + "_" + l.generateRandomString(8)
}

// 生成通知ID
func (l *BaseLogic) generateNotifyId() string {
	return "NOTIFY_" + time.Now().Format("20060102150405") + "_" + l.generateRandomString(8)
}

// 生成随机字符串
func (l *BaseLogic) generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

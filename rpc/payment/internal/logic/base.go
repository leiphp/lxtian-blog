package logic

import (
	"context"
	"lxtian-blog/common/pkg/utils"
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
	return "PAY_" + time.Now().Format("20060102150405") + "_" + utils.GenerateRandomString(8)
}

// 生成退款ID
func (l *BaseLogic) generateRefundId() string {
	return "REF_" + time.Now().Format("20060102150405") + "_" + utils.GenerateRandomString(8)
}

// 生成通知ID
func (l *BaseLogic) generateNotifyId() string {
	return "NOTIFY_" + time.Now().Format("20060102150405") + "_" + utils.GenerateRandomString(8)
}

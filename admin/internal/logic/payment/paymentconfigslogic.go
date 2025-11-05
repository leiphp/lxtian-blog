package payment

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentConfigsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付配置管理
func NewPaymentConfigsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentConfigsLogic {
	return &PaymentConfigsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentConfigsLogic) PaymentConfigs(req *types.PaymentConfigsReq) (resp *types.PaymentConfigsResp, err error) {
	// todo: add your logic here and delete this line

	return
}

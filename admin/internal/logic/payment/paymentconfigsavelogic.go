package payment

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentConfigSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 支付配置保存
func NewPaymentConfigSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentConfigSaveLogic {
	return &PaymentConfigSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentConfigSaveLogic) PaymentConfigSave(req *types.PaymentConfigSaveReq) (resp *types.PaymentConfigSaveResp, err error) {
	// todo: add your logic here and delete this line

	return
}

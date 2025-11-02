package logic

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClosePaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 关闭支付订单
func NewClosePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClosePaymentLogic {
	return &ClosePaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClosePaymentLogic) ClosePayment(req *types.ClosePaymentReq) (resp *types.ClosePaymentResp, err error) {
	// todo: add your logic here and delete this line

	return
}

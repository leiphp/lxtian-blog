package payment

import (
	"context"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryPaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询支付结果
func NewQueryPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPaymentLogic {
	return &QueryPaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryPaymentLogic) QueryPayment(req *types.QueryPaymentReq) (resp *types.QueryPaymentResp, err error) {
	// todo: add your logic here and delete this line

	return
}

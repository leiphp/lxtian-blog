package payment

import (
	"context"

	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResendNotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重发支付通知
func NewResendNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResendNotifyLogic {
	return &ResendNotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResendNotifyLogic) ResendNotify(req *types.ResendNotifyReq) (resp *types.ResendNotifyResp, err error) {
	// todo: add your logic here and delete this line

	return
}

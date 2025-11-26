package web

import (
	"context"
	"fmt"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logc"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderStatLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 订单统计
func NewOrderStatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderStatLogic {
	return &OrderStatLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderStatLogic) OrderStat() (resp *types.OrderStatResp, err error) {
	res, err := l.svcCtx.WebRpc.OrderStat(l.ctx, &web.OrderStatReq{})
	if err != nil {
		logc.Errorf(l.ctx, "OrderStat error message: %s", err)
		return nil, err
	}

	resp = &types.OrderStatResp{
		Data: map[string]interface{}{
			"total_amount": fmt.Sprintf("%.2f", res.GetTotalAmount()),
			"count":        res.GetCount(),
			"month_amount": fmt.Sprintf("%.2f", res.GetMonthAmount()),
		},
	}

	return
}

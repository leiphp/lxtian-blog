package payment

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/payment/pb/payment"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 商品详情
func NewGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodsLogic {
	return &GoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoodsLogic) Goods(req *types.GoodsReq) (resp *types.GoodsResp, err error) {
	res, err := l.svcCtx.PaymentRpc.Goods(l.ctx, &payment.GoodsReq{
		Id: uint64(req.Id),
	})
	if err != nil {
		logc.Errorf(l.ctx, "Goods error: %s", err)
		return nil, err
	}
	resp = new(types.GoodsResp)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	resp.Data = result
	return
}

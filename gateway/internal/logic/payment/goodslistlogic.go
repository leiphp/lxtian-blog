package payment

import (
	"context"
	"encoding/json"
	"lxtian-blog/rpc/payment/pb/payment"

	"github.com/zeromicro/go-zero/core/logc"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GoodsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 商品列表
func NewGoodsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodsListLogic {
	return &GoodsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoodsListLogic) GoodsList(req *types.GoodsListReq) (resp *types.GoodsListResp, err error) {
	// 参数验证和默认值设置
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}

	// 调用 rpc 服务
	res, err := l.svcCtx.PaymentRpc.GoodsList(l.ctx, &payment.GoodsListReq{
		ClassifyId: int32(req.ClassifyId),
		PriceMin:   float32(req.PriceMin),
		PriceMax:   float32(req.PriceMax),
		Page:       uint32(page),
		PageSize:   uint32(pageSize),
		Keywords:   req.Keywords,
		OrderBy:    req.OrderBy,
	})
	if err != nil {
		logc.Errorf(l.ctx, "GoodsList rpc error: %s", err)
		return nil, err
	}

	// 解析返回的 JSON 数据
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		logc.Errorf(l.ctx, "GoodsList json unmarshal error: %s", err)
		return nil, err
	}

	// 构建响应
	resp = new(types.GoodsListResp)
	resp.Page = int(page)
	resp.PageSize = int(pageSize)
	resp.Total = int64(res.GetTotal())
	resp.List = result

	return
}

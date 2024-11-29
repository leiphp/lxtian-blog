package web

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/web/web"
	"time"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 订单列表
func NewOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderListLogic {
	return &OrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderListLogic) OrderList(req *types.OrderListReq) (resp *types.OrderListResp, err error) {
	res, err := l.svcCtx.WebRpc.OrderList(l.ctx, &web.OrderListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		logc.Errorf(l.ctx, "OrderList error message: %s", err)
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		return nil, err
	}
	for k, item := range result {
		if ctimeFloat, ok := item["ctime"].(float64); ok {
			// 将 float64 转为 int64
			ctime := int64(ctimeFloat)
			// 将时间戳转换为 time.Time 类型
			t := time.Unix(ctime, 0)
			// 格式化为 yyyy-mm-dd
			result[k]["date"] = t.Format("2006-01-02")
		}
	}
	resp = new(types.OrderListResp)
	resp.List = result
	resp.Total = uint64(res.GetTotal())
	resp.Page = res.GetPage()
	resp.PageSize = res.GetPageSize()
	return
}

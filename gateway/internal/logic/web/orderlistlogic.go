package web

import (
	"context"
	"encoding/json"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"lxtian-blog/rpc/web/web"
	"time"

	"github.com/zeromicro/go-zero/core/logc"

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
		if createdAtStr, ok := item["created_at"].(string); ok {
			// 解析 ISO 8601 格式的时间字符串，如 "2020-02-03T20:26:35+08:00"
			t, err := time.Parse(time.RFC3339, createdAtStr)
			if err != nil {
				// 如果解析失败，尝试其他常见格式
				t, err = time.Parse("2006-01-02T15:04:05Z07:00", createdAtStr)
				if err != nil {
					l.Errorf("Failed to parse created_at '%s': %v", createdAtStr, err)
					continue
				}
			}
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

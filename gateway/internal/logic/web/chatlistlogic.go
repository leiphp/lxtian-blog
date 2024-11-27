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

type ChatListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 说说列表
func NewChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatListLogic {
	return &ChatListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatListLogic) ChatList(req *types.ChatListReq) (resp *types.ChatListResp, err error) {
	res, err := l.svcCtx.WebRpc.ChatList(l.ctx, &web.ChatListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		logc.Errorf(l.ctx, "ChatList error message: %s", err)
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
			result[k]["date"] = t.Format("01-02")
		}
	}
	resp = new(types.ChatListResp)
	resp.List = result
	resp.Total = uint64(res.GetTotal())
	resp.Page = res.GetPage()
	resp.PageSize = res.GetPageSize()
	return
}

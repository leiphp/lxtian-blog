package web

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/web/web"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 书单详情
func NewBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookLogic {
	return &BookLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BookLogic) Book(req *types.BookReq) (resp *types.BookResp, err error) {
	res, err := l.svcCtx.WebRpc.Book(l.ctx, &web.BookReq{
		Id: req.Id,
	})
	if err != nil {
		logc.Errorf(l.ctx, "Book error: %s", err)
		return nil, err
	}
	resp = new(types.BookResp)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	resp.Data = result
	return
}

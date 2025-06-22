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

type BookListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 书单列表
func NewBookListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookListLogic {
	return &BookListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BookListLogic) BookList(req *types.BookListReq) (resp *types.BookListResp, err error) {
	res, err := l.svcCtx.WebRpc.BookList(l.ctx, &web.BookListReq{
		Column:   req.Column,
		Page:     req.Page,
		PageSize: req.PageSize,
		Keywords: req.Keywords,
	})
	if err != nil {
		logc.Errorf(l.ctx, "ColumnList error message: %s", err)
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		return nil, err
	}
	resp = new(types.BookListResp)
	resp.List = result
	resp.Total = uint64(res.GetTotal())
	resp.Page = res.GetPage()
	resp.PageSize = res.GetPageSize()
	return
}

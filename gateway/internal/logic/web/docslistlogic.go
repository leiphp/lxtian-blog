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

type DocsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 教程列表
func NewDocsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsListLogic {
	return &DocsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsListLogic) DocsList(req *types.DocsListReq) (resp *types.DocsListResp, err error) {
	res, err := l.svcCtx.WebRpc.DocsList(l.ctx, &web.DocsListReq{
		CategoryId: int32(req.CategoryId),
		Level:      req.Level,
		SortBy:     req.SortBy,
		SortOrder:  req.SortOrder,
		Page:       int32(req.Page),
		PageSize:   int32(req.PageSize),
		Keywords:   req.Keywords,
		Tag:        req.Tag,
	})
	if err != nil {
		logc.Errorf(l.ctx, "DocsList error message: %s", err)
		return nil, err
	}
	resp = new(types.DocsListResp)
	var result []*types.DocsItem
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		logc.Errorf(l.ctx, "DocsList unmarshal error: %s", err)
		return nil, err
	}
	resp.List = result
	resp.Total = res.Total
	resp.Page = int(res.Page)
	resp.PageSize = int(res.PageSize)
	return
}

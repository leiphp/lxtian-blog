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

type DocsTagsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 热门标签
func NewDocsTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsTagsLogic {
	return &DocsTagsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsTagsLogic) DocsTags() (resp *types.DocsTagsResp, err error) {
	res, err := l.svcCtx.WebRpc.DocsTags(l.ctx, &web.DocsTagsReq{})
	if err != nil {
		logc.Errorf(l.ctx, "DocsTags error message: %s", err)
		return nil, err
	}
	var result []*types.TagItem
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		return nil, err
	}
	resp = new(types.DocsTagsResp)
	resp.List = result
	return
}

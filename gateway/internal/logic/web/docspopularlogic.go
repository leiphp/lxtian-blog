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

type DocsPopularLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 热门文档
func NewDocsPopularLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsPopularLogic {
	return &DocsPopularLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsPopularLogic) DocsPopular(req *types.DocsPopularReq) (resp *types.DocsPopularResp, err error) {
	res, err := l.svcCtx.WebRpc.DocsPopular(l.ctx, &web.DocsPopularReq{
		Limit: int32(req.Limit),
	})
	if err != nil {
		logc.Errorf(l.ctx, "DocsPopular error message: %s", err)
		return nil, err
	}
	resp = new(types.DocsPopularResp)
	var result []*types.DocsItem
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		logc.Errorf(l.ctx, "DocsPopular unmarshal error: %s", err)
		return nil, err
	}
	resp.List = result
	return
}

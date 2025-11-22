package web

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/web/web"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 教程统计
func NewDocsStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsStatsLogic {
	return &DocsStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsStatsLogic) DocsStats() (resp *types.DocsStatsResp, err error) {
	res, err := l.svcCtx.WebRpc.DocsStats(l.ctx, &web.DocsStatsReq{})
	if err != nil {
		logc.Errorf(l.ctx, "DocsStats error message: %s", err)
		return nil, err
	}
	resp = new(types.DocsStatsResp)
	resp.TotalDocs = res.TotalDocs
	resp.TotalCategories = res.TotalCategories
	resp.TotalViews = res.TotalViews
	resp.TotalLikes = res.TotalLikes
	return
}

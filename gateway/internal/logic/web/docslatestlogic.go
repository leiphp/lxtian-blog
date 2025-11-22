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

type DocsLatestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 最新文档
func NewDocsLatestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsLatestLogic {
	return &DocsLatestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsLatestLogic) DocsLatest(req *types.DocsLatestReq) (resp *types.DocsLatestResp, err error) {
	res, err := l.svcCtx.WebRpc.DocsLatest(l.ctx, &web.DocsLatestReq{
		Limit: int32(req.Limit),
	})
	if err != nil {
		logc.Errorf(l.ctx, "DocsLatest error message: %s", err)
		return nil, err
	}
	resp = new(types.DocsLatestResp)
	var result []*types.DocsItem
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		logc.Errorf(l.ctx, "DocsLatest unmarshal error: %s", err)
		return nil, err
	}
	resp.List = result
	return
}

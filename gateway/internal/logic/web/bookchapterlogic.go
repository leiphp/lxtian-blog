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

type BookChapterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 章节详情
func NewBookChapterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookChapterLogic {
	return &BookChapterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BookChapterLogic) BookChapter(req *types.BookChapterReq) (resp *types.BookChapterResp, err error) {
	res, err := l.svcCtx.WebRpc.BookChapter(l.ctx, &web.BookChapterReq{
		Id: req.Id,
	})
	if err != nil {
		logc.Errorf(l.ctx, "BookChapter error: %s", err)
		return nil, err
	}
	resp = new(types.BookChapterResp)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	resp.Data = result
	return
}

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

type ArticleLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文章喜欢
func NewArticleLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLikeLogic {
	return &ArticleLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleLikeLogic) ArticleLike(req *types.ArticleLikeReq) (resp *types.ArticleLikeResp, err error) {
	res, err := l.svcCtx.WebRpc.ArticleLike(l.ctx, &web.ArticleLikeReq{
		Id: req.Id,
	})
	if err != nil {
		logc.Errorf(l.ctx, "ArticleLike error: %s", err)
		return nil, err
	}
	resp = new(types.ArticleLikeResp)
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	resp.Data = result
	return
}

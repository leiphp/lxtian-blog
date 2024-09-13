package web

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/rpc/web/web"

	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleListLogic {
	return &ArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleListLogic) ArticleList(req *types.ArticleListReq) (resp *types.ArticleListResp, err error) {
	res, err := l.svcCtx.WebRpc.ArticleList(l.ctx, &web.ArticleListReq{
		Cid:      req.Cid,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		logc.Errorf(l.ctx, "ArticleList error message: %s", err)
		return nil, err
	}
	resp = new(types.ArticleListResp)
	// 定义 resp.List 类型
	var respList []map[string]interface{}
	// 遍历 res.List 并转换
	for _, article := range res.List {
		articleMap := map[string]interface{}{
			"id":          article.Id,
			"title":       article.Title,
			"author":      article.Author,
			"description": article.Description,
			"content":     article.Content,
		}
		respList = append(respList, articleMap)
	}
	resp.List = respList
	resp.Total = uint64(res.GetTotal())
	resp.Page = res.GetPage()
	res.PageSize = res.GetPageSize()
	return
}

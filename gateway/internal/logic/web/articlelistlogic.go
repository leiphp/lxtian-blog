package web

import (
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/pkg/utils"
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
	// 结构体转map切片
	respList, err := utils.StructSliceToMapSlice(res.List)
	if err != nil {
		logc.Errorf(l.ctx, "StructSliceToMapSlice: %s", err)
		return nil, err
	}
	resp.List = respList
	resp.Total = uint64(res.GetTotal())
	resp.Page = res.GetPage()
	resp.PageSize = res.GetPageSize()
	return
}

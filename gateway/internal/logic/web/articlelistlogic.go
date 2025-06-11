package web

import (
	"context"
	"encoding/json"
	"fmt"
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
	fmt.Println("req.tid:", req.Tid)
	res, err := l.svcCtx.WebRpc.ArticleList(l.ctx, &web.ArticleListReq{
		Cid:      req.Cid,
		Page:     req.Page,
		PageSize: req.PageSize,
		Types:    req.Types,
		Tid:      req.Tid,
		Keywords: req.Keywords,
	})
	if err != nil {
		logc.Errorf(l.ctx, "ArticleList error message: %s", err)
		return nil, err
	}
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(res.List), &result); err != nil {
		return nil, err
	}
	resp = new(types.ArticleListResp)
	// 结构体转map切片
	//respList, err := utils.StructSliceToMapSlice(res.List)
	//if err != nil {
	//	logc.Errorf(l.ctx, "StructSliceToMapSlice: %s", err)
	//	return nil, err
	//}
	resp.List = result
	resp.Total = uint64(res.GetTotal())
	resp.Page = res.GetPage()
	resp.PageSize = res.GetPageSize()
	return
}

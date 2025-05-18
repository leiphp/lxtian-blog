package web

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文章详情
func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ArticleLogic) Article(req *types.ArticleReq) (resp *types.ArticleResp, err error) {
	res, err := l.svcCtx.WebRpc.Article(l.ctx, &web.ArticleReq{
		Id: req.Id,
	})
	if err != nil {
		logc.Errorf(l.ctx, "Article error: %s", err)
		return nil, err
	}
	resp = new(types.ArticleResp)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(res.Data), &result); err != nil {
		return nil, err
	}
	resp.Data = result
	// 记录订单创建QPS（标签为方法名）
	// svc.OrderCreateQPS.WithLabelValues("CreateOrder").Inc()
	// 记录文章浏览QPS（标签为方法名）
	svc.ArticleViewQPS.WithLabelValues("ArticleView").Inc()
	return
}

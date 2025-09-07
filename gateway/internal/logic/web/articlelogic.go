package web

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"lxtian-blog/rpc/web/web"
	"net/http"

	"github.com/zeromicro/go-zero/core/logc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	req    *http.Request // 添加HTTP请求，用于获取IP
}

// 文章详情
func NewArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ArticleLogic {
	return &ArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// NewArticleLogicWithRequest 创建带HTTP请求的ArticleLogic
func NewArticleLogicWithRequest(ctx context.Context, svcCtx *svc.ServiceContext, req *http.Request) *ArticleLogic {
	return &ArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		req:    req,
	}
}

func (l *ArticleLogic) Article(req *types.ArticleReq) (resp *types.ArticleResp, err error) {
	// 获取客户端IP
	clientIP := ""
	if l.req != nil {
		clientIP = utils.GetClientIP(l.req)
		logc.Infof(l.ctx, "文章 %d 被IP %s 访问", req.Id, clientIP)
	}

	res, err := l.svcCtx.WebRpc.Article(l.ctx, &web.ArticleReq{
		Id:       req.Id,
		ClientIp: clientIP,
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

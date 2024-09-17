package web

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

// 文章详情
func ArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticleReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "ArticleHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := web.NewArticleLogic(r.Context(), svcCtx)
		resp, err := l.Article(&req)
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp.Data, err)
		}
	}
}

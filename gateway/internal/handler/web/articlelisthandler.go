package web

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

func ArticleListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticleListReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "ArticleListHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := web.NewArticleListLogic(r.Context(), svcCtx)
		resp, err := l.ArticleList(&req)
		response.Response(r, w, resp, err)
	}
}

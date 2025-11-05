package content

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
)

// 文章管理
func ArticlesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticlesReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "ArticlesHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := content.NewArticlesLogic(r.Context(), svcCtx)
		resp, err := l.Articles(&req)
		response.Response(r, w, resp, err)
	}
}

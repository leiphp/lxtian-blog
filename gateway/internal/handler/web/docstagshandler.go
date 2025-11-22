package web

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
)

// 热门标签
func DocsTagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := web.NewDocsTagsLogic(r.Context(), svcCtx)
		resp, err := l.DocsTags()
		if err != nil {
			response.Response(r, w, resp, err)
		} else {
			response.Response(r, w, resp.List, err)
		}
	}
}

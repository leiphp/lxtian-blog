package web

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
)

// 教程统计
func DocsStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := web.NewDocsStatsLogic(r.Context(), svcCtx)
		resp, err := l.DocsStats()
		response.Response(r, w, resp, err)
	}
}

package web

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
)

// 专栏列表
func ColumnListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := web.NewColumnListLogic(r.Context(), svcCtx)
		resp, err := l.ColumnList()
		response.Response(r, w, resp.Data, err)
	}
}

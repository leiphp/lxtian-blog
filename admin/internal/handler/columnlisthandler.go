package handler

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/logic"
	"lxtian-blog/admin/internal/svc"
)

// 专栏列表
func ColumnListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewColumnListLogic(r.Context(), svcCtx)
		resp, err := l.ColumnList()
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp.Data, err)
		}
	}
}

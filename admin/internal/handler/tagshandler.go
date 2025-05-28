package handler

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/logic"
	"lxtian-blog/admin/internal/svc"
)

// 标签列表
func TagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewTagsLogic(r.Context(), svcCtx)
		resp, err := l.Tags()
		response.Response(r, w, resp.Data, err)
	}
}

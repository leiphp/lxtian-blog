package web

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
)

// 标签列表
func TagsListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := web.NewTagsListLogic(r.Context(), svcCtx)
		resp, err := l.TagsList()
		response.Response(r, w, resp.Data, err)
	}
}

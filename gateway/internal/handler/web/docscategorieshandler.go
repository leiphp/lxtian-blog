package web

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
)

// 教程分类
func DocsCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := web.NewDocsCategoriesLogic(r.Context(), svcCtx)
		resp, err := l.DocsCategories()
		if err != nil {
			response.Response(r, w, resp, err)
		} else {
			response.Response(r, w, resp.List, err)
		}
	}
}

package content

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
)

// 文章分类
func CategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := content.NewCategoryLogic(r.Context(), svcCtx)
		resp, err := l.Category()
		response.Response(r, w, resp.Data, err)
	}
}

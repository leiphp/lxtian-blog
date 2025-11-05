package content

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
)

// 书单管理
func BooKHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := content.NewBooKLogic(r.Context(), svcCtx)
		resp, err := l.BooK()
		response.Response(r, w, resp.Data, err)
	}
}

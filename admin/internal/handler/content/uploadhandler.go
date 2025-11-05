package content

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
)

// 图片上传
func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := content.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(r)
		response.Response(r, w, resp.Url, err)
	}
}

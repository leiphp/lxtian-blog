package content

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 图片上传
func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadReq
		if err := httpx.ParseForm(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := content.NewUploadLogic(r.Context(), svcCtx)
		resp, err := l.Upload(r, &req)
		if err != nil {
			response.Response(r, w, resp, err)
		} else {
			response.Response(r, w, resp.Url, err)
		}
	}
}

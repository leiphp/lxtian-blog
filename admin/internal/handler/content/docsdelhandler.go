package content

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
)

// 文档删除
func DocsDelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DocsDelReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := content.NewDocsDelLogic(r.Context(), svcCtx)
		resp, err := l.DocsDel(&req)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		response.Response(r, w, resp.Data, nil)
	}
}

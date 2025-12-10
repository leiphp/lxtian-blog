package content

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
)

// 文档列表
func DocsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DocsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := content.NewDocsLogic(r.Context(), svcCtx)
		resp, err := l.Docs(&req)
		response.Response(r, w, resp, err)
	}
}

package web

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

// 教程列表
func DocsListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DocsListReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "DocsListHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := web.NewDocsListLogic(r.Context(), svcCtx)
		resp, err := l.DocsList(&req)
		response.Response(r, w, resp, err)
	}
}

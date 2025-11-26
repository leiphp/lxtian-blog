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

// 文档详情
func DocsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DocsReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "DocsHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := web.NewDocsLogic(r.Context(), svcCtx)
		resp, err := l.Docs(&req, r)
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp.Data, err)
		}
	}
}

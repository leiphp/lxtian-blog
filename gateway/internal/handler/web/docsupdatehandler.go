package web

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

// 文档更新
func DocsUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DocsUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := web.NewDocsUpdateLogic(r.Context(), svcCtx)
		resp, err := l.DocsUpdate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

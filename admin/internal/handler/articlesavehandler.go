package handler

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/logic"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
)

// 文章保存
func ArticleSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ArticleSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "ArticleSaveHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewArticleSaveLogic(r.Context(), svcCtx)
		resp, err := l.ArticleSave(&req)
		response.Response(r, w, resp.Data, err)
	}
}

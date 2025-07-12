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

// 章节保存
func BookChapterSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookChapterSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "BookChapterSaveHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewBookChapterSaveLogic(r.Context(), svcCtx)
		resp, err := l.BookChapterSave(&req)
		response.Response(r, w, resp.Data, err)
	}
}

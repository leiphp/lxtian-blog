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

// 书单详情
func ChapterDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChapterDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "ChapterDetailHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewChapterDetailLogic(r.Context(), svcCtx)
		resp, err := l.ChapterDetail(&req)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		response.Response(r, w, resp.Data, nil)
	}
}

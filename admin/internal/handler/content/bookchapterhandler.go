package content

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
)

// 书单章节
func BookChapterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookChapterReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "BookChapterHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := content.NewBookChapterLogic(r.Context(), svcCtx)
		resp, err := l.BookChapter(&req)
		response.Response(r, w, resp.Data, err)
	}
}

package content

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
)

// 章节数据保存
func BookChapterDataSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BookChapterDataSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := content.NewBookChapterDataSaveLogic(r.Context(), svcCtx)
		resp, err := l.BookChapterDataSave(&req)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		response.Response(r, w, resp.Data, err)
	}
}

package content

import (
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/logic/content"
	"lxtian-blog/admin/internal/svc"
)

// 标签列表
func TagsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TagsReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "TagsHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}
		l := content.NewTagsLogic(r.Context(), svcCtx)
		resp, err := l.Tags(&req)
		response.Response(r, w, resp, err)
	}
}

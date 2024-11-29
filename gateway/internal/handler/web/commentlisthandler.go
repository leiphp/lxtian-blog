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

// 评论列表
func CommentListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CommentListReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "CommentListHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := web.NewCommentListLogic(r.Context(), svcCtx)
		resp, err := l.CommentList(&req)
		response.Response(r, w, resp, err)
	}
}

package user

import (
	"lxtian-blog/admin/internal/logic/user"
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/svc"
)

// 用户信息
func InfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewInfoLogic(r.Context(), svcCtx)
		resp, err := l.Info()
		response.Response(r, w, resp, err)
	}
}

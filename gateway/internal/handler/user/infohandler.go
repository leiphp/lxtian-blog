package user

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/user"
	"lxtian-blog/gateway/internal/svc"
)

func InfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewInfoLogic(r.Context(), svcCtx)
		resp, err := l.Info()
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp, err)
		}
	}
}

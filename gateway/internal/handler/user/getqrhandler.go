package user

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/user"
	"lxtian-blog/gateway/internal/svc"
)

// 获取二维码
func GetqrHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetqrLogic(r.Context(), svcCtx)
		resp, err := l.Getqr()
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp, err)
		}
	}
}

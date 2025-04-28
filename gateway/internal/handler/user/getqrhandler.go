package user

import (
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/common/restful/response"
	"lxtian-blog/gateway/internal/types"
	"net/http"

	"lxtian-blog/gateway/internal/logic/user"
	"lxtian-blog/gateway/internal/svc"
)

// 获取二维码
func GetqrHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetqrReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "GetqrHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}
		l := user.NewGetqrLogic(r.Context(), svcCtx)
		resp, err := l.Getqr(&req)
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp, err)
		}
	}
}

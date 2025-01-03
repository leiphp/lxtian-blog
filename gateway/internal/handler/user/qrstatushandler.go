package user

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/user"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

// 更新扫码状态
func QrStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QrStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewQrStatusLogic(r.Context(), svcCtx)
		resp, err := l.QrStatus(&req)
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp.Data, err)
		}
	}
}

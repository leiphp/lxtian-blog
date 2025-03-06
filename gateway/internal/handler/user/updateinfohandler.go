package user

import (
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/user"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

// 修改用户信息
func UpdateInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "UpdateInfoHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}
		l := user.NewUpdateInfoLogic(r.Context(), svcCtx)
		resp, err := l.UpdateInfo(&req)
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp.Data, err)
		}
	}
}

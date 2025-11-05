package user

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/admin/internal/logic/user"
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
)

// 权限保存
func PremSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PremSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "PremSaveHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := user.NewPremSaveLogic(r.Context(), svcCtx)
		resp, err := l.PremSave(&req)
		response.Response(r, w, resp.Data, err)
	}
}

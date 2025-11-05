package user

import (
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/admin/internal/logic/user"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/svc"
)

// 菜单管理
func MenusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MenusReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "MenusHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}
		l := user.NewMenusLogic(r.Context(), svcCtx)
		resp, err := l.Menus(&req)
		response.Response(r, w, resp.Data, err)
	}
}

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

// 菜单保存
func MenuSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MenuSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "MenuSaveHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := user.NewMenuSaveLogic(r.Context(), svcCtx)
		resp, err := l.MenuSave(&req)
		response.Response(r, w, resp.Data, err)
	}
}

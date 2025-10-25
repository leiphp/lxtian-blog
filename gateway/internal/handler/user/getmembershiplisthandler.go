package user

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/user"
	"lxtian-blog/gateway/internal/svc"
)

// 获取会员列表
func GetMembershipListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewGetMembershipListLogic(r.Context(), svcCtx)
		resp, err := l.GetMembershipList()
		response.Response(r, w, resp, err)
	}
}

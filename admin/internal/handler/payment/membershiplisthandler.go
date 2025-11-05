package payment

import (
	"lxtian-blog/admin/internal/logic/payment"
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/admin/internal/svc"
)

// 会员套餐
func MembershipListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := payment.NewMembershipListLogic(r.Context(), svcCtx)
		resp, err := l.MembershipList()
		response.Response(r, w, resp, err)
	}
}

package payment

import (
	"net/http"

	"lxtian-blog/common/restful/response"
	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RepayOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RepayOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := payment.NewRepayOrderLogic(r.Context(), svcCtx)
		resp, err := l.RepayOrder(&req, r)
		response.Response(r, w, resp, err)
	}
}

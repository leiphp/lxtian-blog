package payment

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

// 关闭订单
func OrdersCloseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrdersCloseReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := payment.NewOrdersCloseLogic(r.Context(), svcCtx)
		resp, err := l.OrdersClose(&req)
		response.Response(r, w, resp, err)
	}
}

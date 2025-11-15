package payment

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

// 取消订单
func OrdersCancelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrdersCancelReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "OrdersCancelHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := payment.NewOrdersCancelLogic(r.Context(), svcCtx)
		resp, err := l.OrdersCancel(&req)
		response.Response(r, w, resp, err)
	}
}

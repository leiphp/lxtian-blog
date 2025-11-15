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

// 删除订单
func OrdersDelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrdersDelReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "OrdersDelHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := payment.NewOrdersDelLogic(r.Context(), svcCtx)
		resp, err := l.OrdersDel(&req)
		response.Response(r, w, resp, err)
	}
}

package payment

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
)

// 支付订单统计
func OrdersStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := payment.NewOrdersStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.OrdersStatistics()
		response.Response(r, w, resp, err)
	}
}

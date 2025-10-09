package payment

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
)

// 支付订单统计
func OrdersStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := payment.NewOrdersStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.OrdersStatistics()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

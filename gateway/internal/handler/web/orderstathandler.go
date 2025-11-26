package web

import (
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/web"
	"lxtian-blog/gateway/internal/svc"
)

// 订单统计
func OrderStatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := web.NewOrderStatLogic(r.Context(), svcCtx)
		resp, err := l.OrderStat()
		if err != nil {
			response.Response(r, w, nil, err)
		} else {
			response.Response(r, w, resp.Data, err)
		}
	}
}

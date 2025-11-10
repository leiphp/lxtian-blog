package payment

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/restful/response"
	"net/http"

	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 支付结果异步通知
func PaymentNotifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PaymentNotifyReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "PaymentNotifyHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := payment.NewPaymentNotifyLogic(r.Context(), svcCtx)
		resp, err := l.PaymentNotify(&req)
		response.Response(r, w, resp, err)
	}
}

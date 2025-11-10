package payment

import (
	"bytes"
	"io"
	"net/http"

	"lxtian-blog/common/restful/response"
	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 支付结果异步通知
func PaymentNotifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PaymentNotifyReq

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logc.Errorf(r.Context(), "PaymentNotifyHandler read body error: %v", err)
			response.Response(r, w, nil, err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "PaymentNotifyHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}
		if len(req.NotifyData) == 0 {
			req.NotifyData = string(bodyBytes)
		}

		l := payment.NewPaymentNotifyLogic(r.Context(), svcCtx)
		resp, err := l.PaymentNotify(&req)
		response.Response(r, w, resp, err)
	}
}

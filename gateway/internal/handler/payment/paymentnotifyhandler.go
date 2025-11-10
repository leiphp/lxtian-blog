package payment

import (
	"bytes"
	"io"
	"net/http"

	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 支付结果异步通知
func PaymentNotifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PaymentNotifyReq

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logx.Errorf("PaymentNotifyHandler read body error: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		logx.Infof("PaymentNotifyHandler request: method=%s url=%s header=%v query=%v body=%s", r.Method, r.URL.String(), r.Header, r.URL.Query(), string(bodyBytes))

		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := payment.NewPaymentNotifyLogic(r.Context(), svcCtx)
		resp, err := l.PaymentNotify(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

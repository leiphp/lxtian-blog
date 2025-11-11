package payment

import (
	"io"
	"lxtian-blog/common/restful/response"
	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"net/http"
	"net/url"

	"github.com/zeromicro/go-zero/core/logc"
)

// 支付结果异步通知
func PaymentNotifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//var req types.PaymentNotifyReq

		// 读取原始body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logc.Errorf(r.Context(), "PaymentNotifyHandler read body error: %v", err)
			response.Response(r, w, nil, err)
			return
		}

		rawBody := string(bodyBytes)
		logc.Infof(r.Context(), "Alipay notify raw body: %s", rawBody)

		// 直接传递原始body作为notify_data，让rpc层处理
		// 同时提取sign和sign_type
		values, _ := url.ParseQuery(rawBody)

		req := types.PaymentNotifyReq{
			NotifyData: rawBody, // 直接传递原始数据
			Sign:       values.Get("sign"),
			SignType:   values.Get("sign_type"),
		}

		logc.Infof(r.Context(), "Passing raw data to rpc, length: %d", len(rawBody))
		logc.Infof(r.Context(), "Sign: %s", req.Sign)
		logc.Infof(r.Context(), "SignType: %s", req.SignType)

		l := payment.NewPaymentNotifyLogic(r.Context(), svcCtx)
		resp, err := l.PaymentNotify(&req)
		response.Response(r, w, resp, err)
	}
}

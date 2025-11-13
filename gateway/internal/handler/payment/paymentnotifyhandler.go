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

		l := payment.NewPaymentNotifyLogic(r.Context(), svcCtx)
		resp, err := l.PaymentNotify(&req)
		if err != nil {
			logc.Errorf(r.Context(), "PaymentNotifyHandler error: %v", err)
			response.Response(r, w, nil, err)
			return
		}

		// 支付宝回调需要返回纯文本 "success"，而不是 JSON
		// 注意：需要先执行 goctl api go 生成 types，然后 resp.Result 才能使用
		result := "success"
		if resp != nil && resp.Result != "" {
			result = resp.Result
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result))
	}
}

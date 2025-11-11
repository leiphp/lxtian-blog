package payment

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"lxtian-blog/common/restful/response"
	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
)

// 支付结果异步通知
func PaymentNotifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req0 types.PaymentNotifyReq
		logc.Infof(r.Context(), "payment notify req: %+v", req0)

		// 读取原始body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logc.Errorf(r.Context(), "PaymentNotifyHandler read body error: %v", err)
			response.Response(r, w, nil, err)
			return
		}

		// 解析为URL参数
		values, err := url.ParseQuery(string(bodyBytes))
		if err != nil {
			logc.Errorf(r.Context(), "PaymentNotifyHandler parse query error: %v", err)
			response.Response(r, w, nil, err)
			return
		}

		// 构建支付宝要求的待验签字符串（排除sign和sign_type，按字母排序）
		var params []string
		for key := range values {
			if key != "sign" && key != "sign_type" {
				// 保持原始值，不要进行URL解码
				params = append(params, fmt.Sprintf("%s=%s", key, values.Get(key)))
			}
		}
		sort.Strings(params)
		notifyData := strings.Join(params, "&")

		// 获取sign和sign_type
		sign := values.Get("sign")
		signType := values.Get("sign_type")

		// 记录调试信息
		logc.Infof(r.Context(), "Alipay notify raw: %s", string(bodyBytes))
		logc.Infof(r.Context(), "Alipay notify data: %s", notifyData)
		logc.Infof(r.Context(), "Alipay sign: %s", sign)
		logc.Infof(r.Context(), "Alipay sign type: %s", signType)

		// 构建请求
		req := types.PaymentNotifyReq{
			NotifyData: notifyData,
			Sign:       sign,
			SignType:   signType,
		}

		l := payment.NewPaymentNotifyLogic(r.Context(), svcCtx)
		resp, err := l.PaymentNotify(&req)
		response.Response(r, w, resp, err)
	}
}

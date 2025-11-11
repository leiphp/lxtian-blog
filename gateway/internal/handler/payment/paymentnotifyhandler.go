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
		var req types.PaymentNotifyReq

		// 读取原始body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logc.Errorf(r.Context(), "PaymentNotifyHandler read body error: %v", err)
			response.Response(r, w, nil, err)
			return
		}

		// 记录原始数据用于调试
		rawBody := string(bodyBytes)
		logc.Infof(r.Context(), "Alipay notify raw body: %s", rawBody)

		// 解析为URL参数
		values, err := url.ParseQuery(rawBody)
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
				// 保持原始值，不要进行任何解码或编码修改
				value := values.Get(key)
				params = append(params, fmt.Sprintf("%s=%s", key, value))
			}
		}
		// 按字母顺序排序（支付宝要求）
		sort.Strings(params)
		notifyData := strings.Join(params, "&")

		// 获取sign和sign_type
		sign := values.Get("sign")
		signType := values.Get("sign_type")

		// URL解码签名（支付宝的签名是Base64编码的，可能包含+号等需要解码）
		decodedSign, err := url.QueryUnescape(sign)
		if err != nil {
			logc.Errorf(r.Context(), "Failed to decode sign, using original: %v", err)
			decodedSign = sign
		}

		// 记录调试信息
		logc.Infof(r.Context(), "Alipay notify data (length: %d): %s", len(notifyData), notifyData)
		logc.Infof(r.Context(), "Alipay sign (original): %s", sign)
		logc.Infof(r.Context(), "Alipay sign (decoded): %s", decodedSign)
		logc.Infof(r.Context(), "Alipay sign type: %s", signType)

		//logc.Infof(r.Context(), "Alipay notify raw: %s", string(bodyBytes))
		//logc.Infof(r.Context(), "Alipay notify data: %s", notifyData)
		//logc.Infof(r.Context(), "Alipay sign: %s", sign)
		//logc.Infof(r.Context(), "Alipay sign type: %s", signType)

		// 构建请求
		req = types.PaymentNotifyReq{
			NotifyData: notifyData,
			Sign:       decodedSign,
			SignType:   signType,
		}

		l := payment.NewPaymentNotifyLogic(r.Context(), svcCtx)
		resp, err := l.PaymentNotify(&req)
		response.Response(r, w, resp, err)
	}
}

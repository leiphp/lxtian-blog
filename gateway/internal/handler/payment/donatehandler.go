package payment

import (
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/restful/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"lxtian-blog/gateway/internal/logic/payment"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
)

// 在线捐赠
func DonateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DonateReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "DonateHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := payment.NewDonateLogic(r.Context(), svcCtx)
		resp, err := l.Donate(&req, r)
		response.Response(r, w, resp, err)
	}
}

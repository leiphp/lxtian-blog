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

// 商品列表
func GoodsListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GoodsListReq
		if err := httpx.Parse(r, &req); err != nil {
			logc.Errorf(r.Context(), "GoodsListHandler error message: %s", err)
			response.Response(r, w, nil, err)
			return
		}

		l := payment.NewGoodsListLogic(r.Context(), svcCtx)
		resp, err := l.GoodsList(&req)
		response.Response(r, w, resp, err)
	}
}

package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

type Body struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Response http返回
func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		//成功返回
		r := &Body{
			Code: 0,
			Msg:  "成功",
			Data: resp,
		}
		httpx.WriteJson(w, http.StatusOK, r)
		return
	}
	//错误返回
	statusCode := http.StatusOK
	errCode := uint32(1)
	// 可以根据错误码，返回具体错误信息
	errMsg := "服务器错误"

	// 判断是否是 HttpError
	if httpErr, ok := err.(*HttpError); ok {
		statusCode = httpErr.StatusCode
		errMsg = httpErr.Message
	} else {
		// 尝试提取 gRPC 错误描述，避免返回 "rpc error: code = Unknown desc = " 前缀
		if st, ok := status.FromError(err); ok {
			errMsg = st.Message()
		} else if err.Error() != "" {
			errMsg = err.Error()
		}
	}

	httpx.WriteJson(w, statusCode, &Body{
		Code: errCode,
		Msg:  errMsg,
		Data: nil,
	})

}

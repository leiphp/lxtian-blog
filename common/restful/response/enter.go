package response

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
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
	} else if err.Error() != "" {
		errMsg = err.Error()
	}

	httpx.WriteJson(w, statusCode, &Body{
		Code: errCode,
		Msg:  errMsg,
		Data: nil,
	})

}

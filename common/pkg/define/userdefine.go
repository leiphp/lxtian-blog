package define

const (
	FailedCode  = 0 //已过期
	DefaultCode = 1 //待扫码
	AlreadyCode = 2 //已扫码
	GoingCode   = 3 //正在登录
	CancelCode  = 4 //取消登录
	SuccessCode = 5 //登录成功
)

var QrLoginCodeMap = map[int]int{
	FailedCode:  0,
	DefaultCode: 1,
	AlreadyCode: 2,
	GoingCode:   3,
	CancelCode:  4,
	SuccessCode: 5,
}

var QrLoginMsgMap = map[int]string{
	FailedCode:  "已过期",
	DefaultCode: "待扫码",
	AlreadyCode: "已扫码",
	GoingCode:   "正在登录",
	CancelCode:  "取消登录",
	SuccessCode: "登录成功",
}

const (
	DefaultLogin = 0 //账号密码登录
	QQLogin      = 1 //QQ登录
	SinaLogin    = 2 //新浪微博登录
	WechatLogin  = 3 //微信扫码登录
	MiniAppLogin = 4 //小程序登录
	GithubLogin  = 5 //GitHub登录
)

// 定义与 JSON 匹配的结构体
type LoginResponse struct {
	Type     string `json:"type"`
	Status   string `json:"status"`
	Msg      string `json:"msg"`
	Token    string `json:"token"`
	UserInfo User   `json:"userInfo"`
}

type User struct {
	Id       int64                  `json:"id"`
	Nickname string                 `json:"nickname"`
	HeadImg  string                 `json:"head_img"`
	Gold     uint64                 `json:"gold"`
	Score    uint64                 `json:"score"`
	Vip      map[string]interface{} `json:"vip"`
}

package userlogic

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/leiphp/wechat/miniapp"
	"github.com/leiphp/wechat/utils"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"lxtian-blog/common/pkg/define"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"
	"strings"
)

type GetqrLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetqrLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetqrLogic {
	return &GetqrLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetqrLogic) Getqr(in *user.GetqrReq) (*user.GetqrResp, error) {
	var (
		page  string
		myApp *miniapp.App
	)
	logc.Info(l.ctx, "MiniAppConf.Appid:", l.svcCtx.Config.MiniAppConf.Appid)
	logc.Info(l.ctx, "MiniAppConf.Secret:", l.svcCtx.Config.MiniAppConf.Secret)
	page = "pages/login/login"
	myApp = miniapp.New(l.svcCtx.Config.MiniAppConf.Appid, l.svcCtx.Config.MiniAppConf.Secret)
	//获取uuid
	uuid := uuid.New().String()
	uuid = strings.Replace(uuid, "-", "", -1)
	//组装二维码参数
	type Acode struct {
		Scene      string      `json:"scene,omitempty"`
		Page       string      `json:"page,omitempty"`
		Width      int         `json:"width,omitempty"`
		AutoColor  bool        `json:"auto_color,omitempty"`
		LineColor  interface{} `json:"line_color,omitempty"`
		IsHyaline  bool        `json:"is_hyaline,omitempty"`
		EnvVersion string      `json:"env_version,omitempty"`
	}

	body := Acode{
		Scene:      uuid,
		Page:       page, //"pages/index/index",
		Width:      280,
		EnvVersion: "release",
		IsHyaline:  true,
	}
	response, err := myApp.PostBody("/wxa/getwxacodeunlimit", utils.JsonToByte(body), true)
	if err != nil {
		return nil, err
	}
	imgBinary := utils.ByteToBase64(response)
	//生成标识

	err = l.svcCtx.Rds.Setex(uuid, fmt.Sprintf(`{"code":%d}`, define.DefaultCode), 5*60)
	if err != nil {
		return nil, err
	}
	//g.Redis().DoVar("SET", shared.ReturnRedisKey(shared.API_CACHE_STRING_QRUUID, uuid), gconv.String("create"), "EX", gconv.String(shared.ApiCacheStringQruuidExpire))

	return &user.GetqrResp{
		Uuid:  uuid,
		QrImg: "data:image/png;base64," + imgBinary,
	}, nil
}

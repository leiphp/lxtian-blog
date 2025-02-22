package userlogic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leiphp/wechat/miniapp"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/define"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/model"
	"lxtian-blog/rpc/user/user"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	switch in.LoginType {
	case define.MiniAppLogin:
		var myApp *miniapp.App
		myApp = miniapp.New(l.svcCtx.Config.MiniAppConf.Appid, l.svcCtx.Config.MiniAppConf.Secret)
		userInfo, err := myApp.Auth().Code2Session(in.Code)
		if err != nil {
			return nil, err
		}
		var txyUser model.TxyUser
		err = l.svcCtx.DB.First(&txyUser, "openid=? and type =?", userInfo.Openid, define.MiniAppLogin).Debug().Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				txyUser.Nickname = ""
				txyUser.HeadImg = ""
				txyUser.Openid = userInfo.Openid
				txyUser.SessionKey = userInfo.Session
				txyUser.MiniappOpenid = userInfo.Appid
				txyUser.Unionid = userInfo.Unionid
				txyUser.Type = 4
				result := l.svcCtx.DB.Create(&txyUser)
				fmt.Println("Id:", result.RowsAffected)
			}
		} else {
			var updateMap = make(map[string]interface{})
			updateMap["session_key"] = userInfo.Session
			if in.Nickname != "" {
				updateMap["nickname"] = in.Nickname
			}
			if in.HeadImg != "" {
				updateMap["head_img"] = in.HeadImg
			}
			l.svcCtx.DB.Model(&txyUser).Where("openid = ? and type = ?", userInfo.Openid, define.MiniAppLogin).Debug().Updates(updateMap)
		}
		err = l.svcCtx.DB.First(&txyUser, "openid = ? and type = ?", userInfo.Openid, define.MiniAppLogin).Debug().Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		res, err := utils.ConvertToLowercaseJSONTags(txyUser)
		if err != nil {
			return nil, err
		}
		jsonData, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}
		return &user.LoginResp{
			Data: string(jsonData),
		}, nil
	default:
		resData, err := l.accountLogin(in)
		if err != nil {
			return nil, err
		}
		return &user.LoginResp{
			Data: *resData,
		}, nil
	}
}

func (l *LoginLogic) accountLogin(in *user.LoginReq) (*string, error) {
	var txyUser model.TxyUser
	err := l.svcCtx.DB.First(&txyUser, "username=?", in.Username).Debug().Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if txyUser.Id == 0 {
		return nil, errors.New("用户名错误")
	}
	// 数据库密码解密
	decodedBytes, err := base64.StdEncoding.DecodeString(txyUser.Password)
	if err != nil {
		return nil, err
	}
	decryptedText, err := utils.Decrypt(decodedBytes)
	if err != nil {
		return nil, err
	}
	if in.Password != decryptedText {
		return nil, errors.New("密码错误！")
	}
	res, err := utils.ConvertToLowercaseJSONTags(txyUser)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	str := string(jsonData)
	return &str, nil
}

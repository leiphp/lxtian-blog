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
	"time"
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
				nickname, err := utils.GenerateNickname(l.svcCtx.DB, 1)
				if err != nil {
					return nil, err
				}
				headImg, err := utils.GetRandomAvatar(1)
				if err != nil {
					return nil, err
				}
				txyUser.Username = "user" + utils.RandomString(8)
				txyUser.Password = l.getPassword("123456")
				txyUser.Nickname = nickname
				txyUser.HeadImg = headImg
				txyUser.Openid = userInfo.Openid
				txyUser.SessionKey = userInfo.Session
				txyUser.MiniappOpenid = userInfo.Appid
				txyUser.Unionid = userInfo.Unionid
				txyUser.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
				txyUser.Type = 4
				result := l.svcCtx.DB.Create(&txyUser)
				fmt.Println("Id:", result.RowsAffected)
			}
		} else {
			var updateMap = make(map[string]interface{})
			updateMap["session_key"] = userInfo.Session
			if in.Nickname != "" && in.Nickname != "微信用户" && in.Nickname != txyUser.Nickname {
				updateMap["nickname"] = in.Nickname
				if in.HeadImg != "" && in.HeadImg != txyUser.HeadImg {
					updateMap["head_img"] = in.HeadImg
				}
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

func (l *LoginLogic) getPassword(password string) string {
	plaintext := []byte(password)
	encryptPassword, err := utils.Encrypt(plaintext)
	if err != nil {
		return ""
	}
	passwordStr := base64.StdEncoding.EncodeToString(encryptPassword)
	return passwordStr
}

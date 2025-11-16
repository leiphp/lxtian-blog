package userlogic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/define"
	"lxtian-blog/common/pkg/oauth"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"
	"time"

	"github.com/leiphp/wechat/miniapp"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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
	case define.QQLogin:
		// QQ登录
		return l.oauthLogin(in, define.QQLogin)
	case define.SinaLogin:
		// 新浪微博登录
		return l.oauthLogin(in, define.SinaLogin)
	case define.WechatLogin:
		// 微信扫码登录
		return l.oauthLogin(in, define.WechatLogin)
	case define.GithubLogin:
		// GitHub登录
		return l.oauthLogin(in, define.GithubLogin)
	case define.MiniAppLogin:
		// 小程序登录
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
				now := time.Now()
				nickname, err := utils.GenerateNickname(l.svcCtx.DB, 1)
				if err != nil {
					return nil, err
				}
				headImg, err := utils.GetRandomAvatar(1)
				if err != nil {
					return nil, err
				}
				username := "user" + utils.RandomString(8)
				txyUser.Username = &username
				txyUser.Password = l.getPassword("123456")
				txyUser.Nickname = nickname
				txyUser.HeadImg = headImg
				txyUser.Openid = userInfo.Openid
				txyUser.SessionKey = userInfo.Session
				txyUser.MiniappOpenid = userInfo.Appid
				txyUser.Unionid = userInfo.Unionid
				txyUser.CreatedAt = &now
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
		// 账号密码登录
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
	if txyUser.ID == 0 {
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

// oauthLogin OAuth社会化登录统一处理
// Gateway已处理OAuth协议交互，RPC只负责用户数据管理
// 接收Gateway传递的用户信息：username(OpenID) + nickname + headImg
func (l *LoginLogic) oauthLogin(in *user.LoginReq, loginType uint32) (*user.LoginResp, error) {
	// 验证必需参数
	if in.Username == "" {
		return nil, fmt.Errorf("OAuth登录必须提供OpenID(username字段)")
	}

	// 构建OAuth用户信息
	oauthUserInfo := &oauth.OAuthUserInfo{
		OpenID:   in.Username, // Username字段存储的是OpenID
		Nickname: in.Nickname,
		HeadImg:  in.HeadImg,
	}

	l.Logger.Infof("OAuth登录 - Type: %d, OpenID: %s, Nickname: %s", loginType, oauthUserInfo.OpenID, oauthUserInfo.Nickname)

	// 查询或创建用户
	var txyUser model.TxyUser
	err := l.svcCtx.DB.First(&txyUser, "openid=? and type=?", oauthUserInfo.OpenID, loginType).Debug().Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新用户
			txyUser, err = l.createOAuthUser(oauthUserInfo, loginType)
			if err != nil {
				return nil, fmt.Errorf("创建用户失败: %w", err)
			}
		} else {
			return nil, err
		}
	} else {
		// 更新用户信息
		err = l.updateOAuthUser(&txyUser, oauthUserInfo)
		if err != nil {
			return nil, fmt.Errorf("更新用户信息失败: %w", err)
		}
	}

	// 重新查询用户信息以获取最新数据
	err = l.svcCtx.DB.First(&txyUser, "openid=? and type=?", oauthUserInfo.OpenID, loginType).Debug().Error
	if err != nil {
		return nil, err
	}

	// 返回用户信息
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
}

// createOAuthUser 创建OAuth登录用户
func (l *LoginLogic) createOAuthUser(oauthUserInfo *oauth.OAuthUserInfo, loginType uint32) (model.TxyUser, error) {
	var txyUser model.TxyUser

	// 生成随机昵称（如果OAuth没有提供）
	nickname := oauthUserInfo.Nickname
	if nickname == "" {
		var err error
		nickname, err = utils.GenerateNickname(l.svcCtx.DB, 1)
		if err != nil {
			return txyUser, err
		}
	}

	// 获取头像（如果OAuth没有提供则使用随机头像）
	headImg := oauthUserInfo.HeadImg
	if headImg == "" {
		var err error
		headImg, err = utils.GetRandomAvatar(1)
		if err != nil {
			return txyUser, err
		}
	}

	// 填充用户信息
	now := time.Now()
	username := "user" + utils.RandomString(8)
	txyUser.Username = &username
	txyUser.Password = l.getPassword("123456") // 默认密码
	txyUser.Nickname = nickname
	txyUser.HeadImg = headImg
	txyUser.Openid = oauthUserInfo.OpenID
	txyUser.AccessToken = oauthUserInfo.AccessToken
	txyUser.Email = oauthUserInfo.Email
	txyUser.Type = int32(loginType)
	txyUser.CreatedAt = &now
	txyUser.Status = 1 // 激活状态

	// 微信特有的UnionID
	if loginType == define.WechatLogin {
		txyUser.Unionid = oauthUserInfo.UnionID
	}

	// 创建用户
	result := l.svcCtx.DB.Create(&txyUser)
	if result.Error != nil {
		return txyUser, result.Error
	}

	l.Logger.Infof("创建OAuth用户成功, OpenID: %s, Type: %d, ID: %d", oauthUserInfo.OpenID, loginType, txyUser.ID)

	return txyUser, nil
}

// updateOAuthUser 更新OAuth登录用户信息
func (l *LoginLogic) updateOAuthUser(txyUser *model.TxyUser, oauthUserInfo *oauth.OAuthUserInfo) error {
	updateMap := make(map[string]interface{})

	// 更新access_token
	if oauthUserInfo.AccessToken != "" {
		updateMap["access_token"] = oauthUserInfo.AccessToken
	}

	// 如果昵称或头像有更新，则更新
	if oauthUserInfo.Nickname != "" && oauthUserInfo.Nickname != txyUser.Nickname {
		updateMap["nickname"] = oauthUserInfo.Nickname
	}

	if oauthUserInfo.HeadImg != "" && oauthUserInfo.HeadImg != txyUser.HeadImg {
		updateMap["head_img"] = oauthUserInfo.HeadImg
	}

	// 更新邮箱
	if oauthUserInfo.Email != "" && oauthUserInfo.Email != txyUser.Email {
		updateMap["email"] = oauthUserInfo.Email
	}

	// 更新UnionID（微信）
	if oauthUserInfo.UnionID != "" && oauthUserInfo.UnionID != txyUser.Unionid {
		updateMap["unionid"] = oauthUserInfo.UnionID
	}

	// 更新最后登录时间和次数
	updateMap["last_login_time"] = time.Now().Unix()
	updateMap["login_times"] = txyUser.LoginTimes + 1

	// 执行更新
	if len(updateMap) > 0 {
		err := l.svcCtx.DB.Model(txyUser).Where("openid=? and type=?", oauthUserInfo.OpenID, txyUser.Type).Updates(updateMap).Error
		if err != nil {
			return err
		}
		l.Logger.Infof("更新OAuth用户信息成功, OpenID: %s, Type: %d", oauthUserInfo.OpenID, txyUser.Type)
	}

	return nil
}

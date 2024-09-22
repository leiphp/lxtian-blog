package userlogic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	encrypt "lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/model"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	if in.Code != "000000" {
		return nil, errors.New("验证码校验失败！")
	}
	// 查询用户名
	var txyUser model.TxyUser
	err := l.svcCtx.DB.First(&txyUser, "username = ?", in.Username).Debug().Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if txyUser.Id > 0 {
		return nil, errors.New("用户名已存在！")
	}
	// 插入数据
	plaintext := []byte(in.Password)
	encryptPassword, err := encrypt.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	passwordStr := base64.StdEncoding.EncodeToString(encryptPassword)
	userRes := model.TxyUser{
		Username: in.Username,
		Password: passwordStr,
	}
	res := l.svcCtx.DB.Create(&userRes)
	if res.Error != nil {
		return nil, res.Error
	}
	var userInfo = make(map[string]interface{}, 3)
	userInfo["id"] = userRes.Id
	userInfo["username"] = userRes.Username
	userInfo["nickname"] = userRes.Nickname
	// json格式化
	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		return nil, err
	}
	return &user.RegisterResp{
		Data: string(jsonData),
	}, nil
}

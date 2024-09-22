package userlogic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/user/model"

	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
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
	return &user.LoginResp{
		Data: string(jsonData),
	}, nil
}

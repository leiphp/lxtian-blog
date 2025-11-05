package user

import (
	"context"
	"encoding/base64"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/jwts"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 后台登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	//var txyUser mysql.TxyUser
	var result struct {
		mysql.TxyUser
		Key         string `json:"key"`
		Permissions string `json:"permissions"`
	}
	err = l.svcCtx.DB.
		Model(&mysql.TxyUser{}).
		Select("txy_user.id,nickname,username,password,is_admin,head_img,type,r.key,GROUP_CONCAT(rp.perm_id) AS permissions").
		Joins("left join txy_user_roles as ur on ur.user_id = txy_user.id").
		Joins("left join txy_roles as r on r.id = ur.role_id").
		Joins("left join txy_role_permissions  as rp on rp.role_id = r.id").
		Where("username = ?", req.Username).
		Debug().
		First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名不存在")
		}
		return nil, err // 其他数据库错误
	}

	permList := strings.Split(result.Permissions, ",")

	// 数据库密码解密
	decodedBytes, err := base64.StdEncoding.DecodeString(result.Password)
	if err != nil {
		return nil, err
	}
	decryptedText, err := utils.Decrypt(decodedBytes)
	if err != nil {
		return nil, err
	}
	if req.Password != decryptedText {
		return nil, errors.New("密码错误！")
	}

	// 获取token
	auth := l.svcCtx.Config.Auth
	token, err := jwts.GenToken(jwts.JwtPayLoad{
		UserID:   uint(result.Id),
		Username: result.Username,
	}, auth.AccessSecret, auth.AccessExpire)
	if err != nil {
		return nil, err
	}

	resp = new(types.LoginResp)
	resp.Token = token
	resp.User = types.User{
		Id:          int(result.Id),
		Username:    result.Username,
		Role:        result.Key,
		Permissions: permList,
	}
	return
}

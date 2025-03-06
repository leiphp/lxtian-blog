package userlogic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"lxtian-blog/rpc/user/model"

	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInfoLogic {
	return &UpdateInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateInfoLogic) UpdateInfo(in *user.UpdateInfoReq) (*user.UpdateInfoResp, error) {
	fmt.Println("aaaaaaaa:", in)
	// 查询用户名
	var txyUser model.TxyUser
	err := l.svcCtx.DB.Debug().First(&txyUser, "id = ?", in.Id).Error
	fmt.Println("err33333:", err)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	fmt.Println("txyUser3333:", txyUser)
	if txyUser.Id == 0 {
		return nil, errors.New("用户不存在！")
	}
	// 更新数据
	userRes := model.TxyUser{
		Nickname: in.Nickname,
		HeadImg:  in.HeadImg,
	}
	fmt.Println("userRes:", userRes)
	res := l.svcCtx.DB.Debug().Where("id = ?", in.Id).Updates(userRes)
	if res.Error != nil {
		return nil, res.Error
	}
	var userInfo = make(map[string]interface{}, 3)
	userInfo["id"] = txyUser.Id
	// json格式化
	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		return nil, err
	}
	return &user.UpdateInfoResp{
		Data: string(jsonData),
	}, nil
}

package userlogic

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"lxtian-blog/common/model"

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
	// 查询用户名
	var txyUser model.TxyUser
	err := l.svcCtx.DB.Debug().First(&txyUser, "id = ?", in.Id).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if txyUser.ID == 0 {
		return nil, errors.New("用户不存在！")
	}
	// 更新数据
	userRes := model.TxyUser{
		Nickname: in.Nickname,
		HeadImg:  in.HeadImg,
	}
	res := l.svcCtx.DB.Debug().Where("id = ?", in.Id).Updates(userRes)
	if res.Error != nil {
		return nil, res.Error
	}
	var userInfo = make(map[string]interface{}, 3)
	userInfo["id"] = txyUser.ID
	// json格式化
	jsonData, err := json.Marshal(userInfo)
	if err != nil {
		return nil, err
	}
	return &user.UpdateInfoResp{
		Data: string(jsonData),
	}, nil
}

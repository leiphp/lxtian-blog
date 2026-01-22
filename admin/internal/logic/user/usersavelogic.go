package user

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户保存
func NewUserSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSaveLogic {
	return &UserSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSaveLogic) UserSave(req *types.UserSaveReq) (resp *types.UserSaveResp, err error) {
	// id 必须大于 0，因为是自增主键
	if req.Id <= 0 {
		return nil, errors.New("用户ID必须大于0")
	}

	// 先查询用户是否存在
	var user mysql.TxyUser
	err = l.svcCtx.DB.Where("id = ?", req.Id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在: id=%d", req.Id)
		}
		l.Errorf("查询用户失败: id=%d, err=%v", req.Id, err)
		return nil, err
	}

	// 构建更新字段，只包含非默认值的字段
	updates := make(map[string]interface{})

	// 昵称（必填字段，总是更新）
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}

	// 密码：如果提供了且不为空，则加密并更新
	if req.Password != "" {
		encryptedBytes, err := utils.Encrypt([]byte(req.Password))
		if err != nil {
			l.Errorf("密码加密失败: err=%v", err)
			return nil, fmt.Errorf("密码加密失败: %w", err)
		}
		updates["password"] = base64.StdEncoding.EncodeToString(encryptedBytes)
	}

	// 邮箱：如果提供了且不为空，则更新
	if req.Email != "" {
		updates["email"] = req.Email
	}

	// 状态：如果提供了且大于0，则更新（0 可能是禁用状态，但根据要求，0 视为默认值不更新）
	if req.Status > 0 {
		updates["status"] = uint64(req.Status)
	}

	// 金币：如果提供了且大于0，则更新
	if req.Gold > 0 {
		updates["gold"] = uint64(req.Gold)
	}

	// 积分：如果提供了且大于0，则更新
	if req.Score > 0 {
		updates["score"] = uint64(req.Score)
	}

	// 更新时间
	updates["updated_at"] = time.Now()

	// 更新用户表
	if len(updates) > 0 {
		err = l.svcCtx.DB.Model(&mysql.TxyUser{}).
			Where("id = ?", req.Id).
			Updates(updates).Error
		if err != nil {
			l.Errorf("更新用户失败: id=%d, err=%v", req.Id, err)
			return nil, err
		}
		l.Infof("更新用户成功: id=%d", req.Id)
	}

	// 角色：如果提供了且大于0，则更新用户角色关联表
	if req.RoleId > 0 {
		// 先删除旧的角色关联
		err = l.svcCtx.DB.Table("txy_user_roles").
			Where("user_id = ?", req.Id).
			Delete(nil).Error
		if err != nil {
			l.Errorf("删除用户角色关联失败: user_id=%d, err=%v", req.Id, err)
			return nil, err
		}

		// 插入新的角色关联
		err = l.svcCtx.DB.Table("txy_user_roles").
			Create(map[string]interface{}{
				"user_id": req.Id,
				"role_id": req.RoleId,
			}).Error
		if err != nil {
			l.Errorf("创建用户角色关联失败: user_id=%d, role_id=%d, err=%v", req.Id, req.RoleId, err)
			return nil, err
		}
		l.Infof("更新用户角色关联成功: user_id=%d, role_id=%d", req.Id, req.RoleId)
	}

	resp = new(types.UserSaveResp)
	resp.Data = true
	return
}

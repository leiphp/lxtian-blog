package userlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/common/repository/user_repo"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoLogic {
	return &InfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *InfoLogic) Info(in *user.InfoReq) (*user.InfoResp, error) {
	// 使用全局 ServiceContext 中的本地缓存（30 分钟过期）
	cache := l.svcCtx.Cache
	// 查询并处理会员信息（优先从 Redis 获取，未命中再查 DB）
	membershipRepo := user_repo.NewUserMembershipRepository(l.svcCtx.DB, l.svcCtx.Rds)
	membershipInfo, err := membershipRepo.GetActiveMembershipByUserId(l.ctx, int64(in.Id))
	if err != nil {
		l.Errorf("Failed to get membership info: %v", err)
		// 会员信息查询失败不影响用户信息返回，继续执行
		membershipInfo = nil
	}

	// 尝试从缓存获取用户信息
	v, exist := cache.Get(fmt.Sprintf("userInfo:%d", in.Id))
	if exist {
		l.Infof("从缓存中获取用户信息: %v", v)
		// 缓存中存在数据，合并最新的会员信息后返回
		res := l.mergeUserAndMembership(v, membershipInfo)
		return l.returnData(res)
	}

	// 如果缓存中没有，查询数据库
	txyUser, err := l.getUserFromDB(in.Id)
	if err != nil {
		return nil, err
	}
	// 转换为JSON小写标签格式
	userData, err := utils.ConvertToLowercaseJSONTags(txyUser)
	if err != nil {
		return nil, err
	}
	l.Infof("从DB中获取用户信息: %v", userData)
	// 合并用户信息和会员信息
	res := l.mergeUserAndMembership(userData, membershipInfo)

	// 设置缓存
	cache.Set(fmt.Sprintf("userInfo:%d", in.Id), res)
	// 返回结果
	return l.returnData(res)
}

// 从数据库获取用户信息
func (l *InfoLogic) getUserFromDB(id uint32) (*model.TxyUser, error) {
	var txyUser model.TxyUser
	if err := l.svcCtx.DB.
		Select("id", "username", "nickname", "head_img", "email", "gold", "score", "type", "status").
		Where("id = ?", id).
		First(&txyUser).Error; err != nil {
		return nil, err
	}
	return &txyUser, nil
}

// 查询出表的全部字段
func (l *InfoLogic) getAllUserFromDB(id uint32) (*model.TxyUser, error) {
	var txyUser model.TxyUser
	if err := l.svcCtx.DB.First(&txyUser, "id = ?", id).Debug().Error; err != nil {
		return nil, err
	}
	return &txyUser, nil
}

// 合并用户信息和会员信息
func (l *InfoLogic) mergeUserAndMembership(userData interface{}, membershipInfo map[string]interface{}) map[string]interface{} {
	// 将用户数据转换为 map
	userBytes, _ := json.Marshal(userData)
	var userMap map[string]interface{}
	json.Unmarshal(userBytes, &userMap)

	// 如果会员信息存在，添加到用户信息中
	if membershipInfo != nil {
		userMap["membership"] = membershipInfo
	} else {
		// 如果没有会员信息，返回 null
		userMap["membership"] = nil
	}

	return userMap
}

// 返回InfoResp
func (l *InfoLogic) returnData(data interface{}) (*user.InfoResp, error) {
	// 转换缓存中的数据为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return &user.InfoResp{
		Data: string(jsonData),
	}, nil
}

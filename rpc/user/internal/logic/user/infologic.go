package userlogic

import (
	"context"
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
	cacheKey := fmt.Sprintf("userInfo:%d", in.Id)
	v, exist := cache.Get(cacheKey)
	if exist {
		l.Infof("从缓存中获取用户信息: %v", v)
		// 缓存中存在数据，转换为 UserInfo 并合并最新的会员信息后返回
		userInfo, err := l.convertCacheToUserInfo(v)
		if err != nil {
			l.Errorf("Failed to convert cache data to UserInfo: %v", err)
			// 转换失败，继续从数据库获取
		} else {
			// 构建会员信息
			membershipInfoProto := l.buildMembershipInfo(membershipInfo)
			return &user.InfoResp{
				User:       userInfo,
				Membership: membershipInfoProto,
			}, nil
		}
	}

	// 如果缓存中没有，查询数据库
	txyUser, err := l.getUserFromDB(in.Id)
	if err != nil {
		return nil, err
	}

	// 构建用户信息
	userInfo := l.buildUserInfo(txyUser)

	// 设置缓存
	cache.Set(cacheKey, txyUser)

	// 构建会员信息
	membershipInfoProto := l.buildMembershipInfo(membershipInfo)

	// 返回结构化的响应
	return &user.InfoResp{
		User:       userInfo,
		Membership: membershipInfoProto,
	}, nil
}

// 从数据库获取用户信息
func (l *InfoLogic) getUserFromDB(id uint32) (*model.TxyUser, error) {
	var txyUser model.TxyUser
	if err := l.svcCtx.DB.
		Select("id", "uid", "username", "nickname", "head_img", "email", "gold", "score", "type", "status").
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

// buildUserInfo 从 TxyUser 构建 UserInfo
func (l *InfoLogic) buildUserInfo(txyUser *model.TxyUser) *user.UserInfo {
	username := ""
	if txyUser.Username != nil {
		username = *txyUser.Username
	}
	return &user.UserInfo{
		Id:       uint64(txyUser.ID),
		Uid:      uint64(txyUser.UID),
		Username: username,
		Nickname: txyUser.Nickname,
		Email:    txyUser.Email,
		HeadImg:  txyUser.HeadImg,
		Gold:     uint64(txyUser.Gold),
		Score:    uint64(txyUser.Score),
		Type:     uint64(txyUser.Type),
		Status:   uint64(txyUser.Status),
	}
}

// convertCacheToUserInfo 将缓存中的数据转换为 UserInfo
func (l *InfoLogic) convertCacheToUserInfo(v interface{}) (*user.UserInfo, error) {
	// 尝试类型断言为 TxyUser
	if txyUser, ok := v.(*model.TxyUser); ok {
		return l.buildUserInfo(txyUser), nil
	}
	// 如果不是 TxyUser 类型，返回错误
	return nil, fmt.Errorf("cache data is not *model.TxyUser type")
}

// buildMembershipInfo 构建会员信息
func (l *InfoLogic) buildMembershipInfo(membershipInfo map[string]interface{}) *user.MembershipInfo {
	if membershipInfo == nil {
		return nil
	}
	return &user.MembershipInfo{
		IsValid:   utils.GetBoolValue(membershipInfo, "is_valid"),
		IsActive:  utils.GetInt32Value(membershipInfo, "is_active"),
		Level:     utils.GetInt32Value(membershipInfo, "level"),
		StartTime: utils.GetStringValue(membershipInfo, "start_time"),
		EndTime:   utils.GetStringValue(membershipInfo, "end_time"),
		TypeId:    utils.GetInt64Value(membershipInfo, "type_id"),
		TotalDays: utils.GetInt32Value(membershipInfo, "total_days"),
	}
}

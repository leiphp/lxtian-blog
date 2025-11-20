package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/leiphp/unit-go-sdk/pkg/gconv"
	"github.com/zeromicro/go-zero/core/logc"
	"lxtian-blog/common/pkg/define"
	"lxtian-blog/common/pkg/jwts"
	"lxtian-blog/common/pkg/redis"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	var res *user.LoginResp
	var token string
	var message string
	switch req.LoginType {
	case define.MiniAppLogin: //小程序登录
		if req.Uuid != "" {
			jsonStr := fmt.Sprintf(`{"code":%d,"token":"%s"}`, define.GoingCode, token)
			l.svcCtx.Rds.Setex(redis.ReturnRedisKey(redis.UserScanString, req.Uuid), jsonStr, 5*60)
			message, err = utils.GetSocketMessage(token, "正在登录", define.User{})
			if err != nil {
				return nil, err
			}
			// 根据二维码id获取ws_user_id
			wsUserId, err := l.svcCtx.Rds.Get(redis.ReturnRedisKey(redis.WsUserIdString, req.Uuid))
			if err != nil {
				return nil, err
			}
			utils.SendMessageToChatService(l.svcCtx.Config.WsService.Host, l.svcCtx.Config.WsService.Port, wsUserId, message)
		}
		res, err = l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
			LoginType: uint32(req.LoginType),
			Code:      req.Code,
			Nickname:  gconv.String(req.Userinfo["nickName"]),
			HeadImg:   gconv.String(req.Userinfo["avatarUrl"]),
		})
		fmt.Println("res:", res)
		if err != nil {
			logc.Errorf(l.ctx, "Login error message: %s", err)
			return nil, err
		}
		var result map[string]interface{}
		if err = json.Unmarshal([]byte(res.Data), &result); err != nil {
			return nil, err
		}
		fmt.Println("result:", result)
		// 获取token
		auth := l.svcCtx.Config.Auth
		token, err = jwts.GenToken(jwts.JwtPayLoad{
			UserID:   uint(result["id"].(float64)),
			Username: result["username"].(string),
		}, auth.AccessSecret, auth.AccessExpire)
		if err != nil {
			return nil, err
		}

		err = l.svcCtx.Rds.Setex(redis.ReturnRedisKey(redis.UserTokenString, result["id"]), token, int(auth.AccessExpire)*3600)
		if err != nil {
			return nil, err
		}

		// 如果是扫码登录，则设置uuid的状态
		if req.Uuid != "" {
			jsonStr := fmt.Sprintf(`{"code":%d,"token":"%s"}`, define.SuccessCode, token)
			l.svcCtx.Rds.Setex(redis.ReturnRedisKey(redis.UserScanString, req.Uuid), jsonStr, 5*60)
			userinfo := define.User{
				Id:       gconv.Int64(result["id"]),
				Nickname: gconv.String(result["nickname"]),
				HeadImg:  gconv.String(result["head_img"]),
				Gold:     gconv.Uint64(result["gold"]),
				Score:    gconv.Uint64(result["score"]),
				Vip:      gconv.Map(result["vip"]),
			}
			message, err = utils.GetSocketMessage(token, "登录成功", userinfo)
			if err != nil {
				return nil, err
			}
			// 根据二维码id获取ws_user_id
			wsUserId, err := l.svcCtx.Rds.Get(redis.ReturnRedisKey(redis.WsUserIdString, req.Uuid))
			if err != nil {
				return nil, err
			}
			utils.SendMessageToChatService(l.svcCtx.Config.WsService.Host, l.svcCtx.Config.WsService.Port, wsUserId, message)
		}

		resp = new(types.LoginResp)
		resp.User = result
		resp.AccessToken = token
		resp.ExpiresIn = uint64(auth.AccessExpire) * 3600
		return
	default: //账号登录
		res, err = l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			logc.Errorf(l.ctx, "Login error message: %s", err)
			return nil, err
		}
		var result map[string]interface{}
		if err = json.Unmarshal([]byte(res.Data), &result); err != nil {
			return nil, err
		}
		// 获取token
		auth := l.svcCtx.Config.Auth
		token, err = jwts.GenToken(jwts.JwtPayLoad{
			UserID:   uint(result["id"].(float64)),
			Username: result["username"].(string),
		}, auth.AccessSecret, auth.AccessExpire)
		if err != nil {
			return nil, err
		}
		resp = new(types.LoginResp)
		resp.User = result
		resp.AccessToken = token
		resp.ExpiresIn = uint64(auth.AccessExpire) * 3600
		return
	}
	return &types.LoginResp{}, nil
}

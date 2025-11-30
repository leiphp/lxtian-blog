package web

import (
	"context"
	"errors"
	"lxtian-blog/gateway/internal/svc"
	"lxtian-blog/gateway/internal/types"
	"lxtian-blog/rpc/user/user"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DocsUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文档更新
func NewDocsUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsUpdateLogic {
	return &DocsUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsUpdateLogic) DocsUpdate(req *types.DocsUpdateReq) (resp *types.DocsUpdateResp, err error) {
	// 1. 从中间件获取用户ID
	userId, ok := l.ctx.Value("user_id").(uint)
	if !ok {
		logx.Errorf("DocsUpdate user_id not found in context")
		return nil, errors.New("请先登录")
	}

	// 2. 通过 UserRpc.Info 查询用户信息
	userInfo, err := l.svcCtx.UserRpc.Info(l.ctx, &user.InfoReq{
		Id: uint32(userId),
	})
	if err != nil {
		logc.Errorf(l.ctx, "DocsUpdate get user info error: %s", err)
		return nil, errors.New("获取用户信息失败")
	}

	// 3. 检查用户信息是否存在
	if userInfo == nil || userInfo.User == nil {
		logc.Errorf(l.ctx, "DocsUpdate user info is nil")
		return nil, errors.New("用户信息不存在")
	}

	// 4. 判断用户是否是管理员（is_admin==1 时，Role 会被设置为 "administrator"）
	if userInfo.User.Role != "administrator" {
		logc.Errorf(l.ctx, "DocsUpdate user %d is not admin, role: %s", userId, userInfo.User.Role)
		return nil, errors.New("无权限修改文档，仅管理员可操作")
	}

	logx.Infof("用户 %d (管理员) 正在更新文档 %d", userId, req.Id)

	// 5. 权限验证通过，调用 WebRpc 执行文档更新
	updateResp, err := l.svcCtx.WebRpc.DocsUpdate(l.ctx, &web.DocsUpdateReq{
		Id:      req.Id,
		Content: req.Content,
	})
	if err != nil {
		logc.Errorf(l.ctx, "DocsUpdate call WebRpc error: %s", err)
		return nil, errors.New("更新文档失败")
	}

	resp = &types.DocsUpdateResp{
		Status: updateResp.Status,
		Id:     updateResp.Id,
	}

	return resp, nil
}

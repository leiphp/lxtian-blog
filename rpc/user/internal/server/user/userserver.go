// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2
// Source: user.proto

package server

import (
	"context"

	"lxtian-blog/rpc/user/internal/logic/user"
	"lxtian-blog/rpc/user/internal/svc"
	"lxtian-blog/rpc/user/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

func (s *UserServer) Getqr(ctx context.Context, in *user.GetqrReq) (*user.GetqrResp, error) {
	l := userlogic.NewGetqrLogic(ctx, s.svcCtx)
	return l.Getqr(in)
}

func (s *UserServer) QrStatus(ctx context.Context, in *user.QrStatusReq) (*user.QrStatusResp, error) {
	l := userlogic.NewQrStatusLogic(ctx, s.svcCtx)
	return l.QrStatus(in)
}

func (s *UserServer) Register(ctx context.Context, in *user.RegisterReq) (*user.RegisterResp, error) {
	l := userlogic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

func (s *UserServer) Login(ctx context.Context, in *user.LoginReq) (*user.LoginResp, error) {
	l := userlogic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *UserServer) Info(ctx context.Context, in *user.InfoReq) (*user.InfoResp, error) {
	l := userlogic.NewInfoLogic(ctx, s.svcCtx)
	return l.Info(in)
}

func (s *UserServer) UpdateInfo(ctx context.Context, in *user.UpdateInfoReq) (*user.UpdateInfoResp, error) {
	l := userlogic.NewUpdateInfoLogic(ctx, s.svcCtx)
	return l.UpdateInfo(in)
}

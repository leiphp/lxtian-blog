// Code generated by goctl. DO NOT EDIT.
// Source: user.proto

package user

import (
	"context"

	"lxtian-blog/rpc/user/user"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	InfoReq      = user.InfoReq
	InfoResp     = user.InfoResp
	LoginReq     = user.LoginReq
	LoginResp    = user.LoginResp
	RegisterReq  = user.RegisterReq
	RegisterResp = user.RegisterResp

	User interface {
		Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error)
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
		Info(ctx context.Context, in *InfoReq, opts ...grpc.CallOption) (*InfoResp, error)
	}

	defaultUser struct {
		cli zrpc.Client
	}
)

func NewUser(cli zrpc.Client) User {
	return &defaultUser{
		cli: cli,
	}
}

func (m *defaultUser) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Register(ctx, in, opts...)
}

func (m *defaultUser) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultUser) Info(ctx context.Context, in *InfoReq, opts ...grpc.CallOption) (*InfoResp, error) {
	client := user.NewUserClient(m.cli.Conn())
	return client.Info(ctx, in, opts...)
}

// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2
// Source: web.proto

package web

import (
	"context"

	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	ArticleListReq  = web.ArticleListReq
	ArticleListResp = web.ArticleListResp
	ArticleReq      = web.ArticleReq
	ArticleResp     = web.ArticleResp

	Web interface {
		ArticleList(ctx context.Context, in *ArticleListReq, opts ...grpc.CallOption) (*ArticleListResp, error)
		Article(ctx context.Context, in *ArticleReq, opts ...grpc.CallOption) (*ArticleResp, error)
	}

	defaultWeb struct {
		cli zrpc.Client
	}
)

func NewWeb(cli zrpc.Client) Web {
	return &defaultWeb{
		cli: cli,
	}
}

func (m *defaultWeb) ArticleList(ctx context.Context, in *ArticleListReq, opts ...grpc.CallOption) (*ArticleListResp, error) {
	client := web.NewWebClient(m.cli.Conn())
	return client.ArticleList(ctx, in, opts...)
}

func (m *defaultWeb) Article(ctx context.Context, in *ArticleReq, opts ...grpc.CallOption) (*ArticleResp, error) {
	client := web.NewWebClient(m.cli.Conn())
	return client.Article(ctx, in, opts...)
}

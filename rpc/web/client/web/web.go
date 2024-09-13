// Code generated by goctl. DO NOT EDIT.
// Source: web.proto

package web

import (
	"context"

	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	Article         = web.Article
	ArticleListReq  = web.ArticleListReq
	ArticleListResp = web.ArticleListResp

	Web interface {
		ArticleList(ctx context.Context, in *ArticleListReq, opts ...grpc.CallOption) (*ArticleListResp, error)
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

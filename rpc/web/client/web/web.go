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
	ArticleListReq   = web.ArticleListReq
	ArticleListResp  = web.ArticleListResp
	ArticleReq       = web.ArticleReq
	ArticleResp      = web.ArticleResp
	CategoryListReq  = web.CategoryListReq
	CategoryListResp = web.CategoryListResp
	ChatListReq      = web.ChatListReq
	ChatListResp     = web.ChatListResp
	CommentListReq   = web.CommentListReq
	CommentListResp  = web.CommentListResp
	OrderListReq     = web.OrderListReq
	OrderListResp    = web.OrderListResp
	TagsListReq      = web.TagsListReq
	TagsListResp     = web.TagsListResp

	Web interface {
		ArticleList(ctx context.Context, in *ArticleListReq, opts ...grpc.CallOption) (*ArticleListResp, error)
		Article(ctx context.Context, in *ArticleReq, opts ...grpc.CallOption) (*ArticleResp, error)
		CategoryList(ctx context.Context, in *CategoryListReq, opts ...grpc.CallOption) (*CategoryListResp, error)
		ChatList(ctx context.Context, in *ChatListReq, opts ...grpc.CallOption) (*ChatListResp, error)
		CommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error)
		OrderList(ctx context.Context, in *OrderListReq, opts ...grpc.CallOption) (*OrderListResp, error)
		TagsList(ctx context.Context, in *TagsListReq, opts ...grpc.CallOption) (*TagsListResp, error)
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

func (m *defaultWeb) CategoryList(ctx context.Context, in *CategoryListReq, opts ...grpc.CallOption) (*CategoryListResp, error) {
	client := web.NewWebClient(m.cli.Conn())
	return client.CategoryList(ctx, in, opts...)
}

func (m *defaultWeb) ChatList(ctx context.Context, in *ChatListReq, opts ...grpc.CallOption) (*ChatListResp, error) {
	client := web.NewWebClient(m.cli.Conn())
	return client.ChatList(ctx, in, opts...)
}

func (m *defaultWeb) CommentList(ctx context.Context, in *CommentListReq, opts ...grpc.CallOption) (*CommentListResp, error) {
	client := web.NewWebClient(m.cli.Conn())
	return client.CommentList(ctx, in, opts...)
}

func (m *defaultWeb) OrderList(ctx context.Context, in *OrderListReq, opts ...grpc.CallOption) (*OrderListResp, error) {
	client := web.NewWebClient(m.cli.Conn())
	return client.OrderList(ctx, in, opts...)
}

func (m *defaultWeb) TagsList(ctx context.Context, in *TagsListReq, opts ...grpc.CallOption) (*TagsListResp, error) {
	client := web.NewWebClient(m.cli.Conn())
	return client.TagsList(ctx, in, opts...)
}

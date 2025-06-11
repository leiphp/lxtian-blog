package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommentListLogic {
	return &CommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommentListLogic) CommentList(in *web.CommentListReq) (*web.CommentListResp, error) {
	where := map[string]interface{}{}
	where["status"] = 1
	where["pid"] = 0
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	offset := (in.Page - 1) * in.PageSize
	var results []map[string]interface{}
	err := l.svcCtx.DB.
		Table("txy_comment as c").
		Select("c.id,c.ouid,c.aid,c.ctime,c.content,c.status,u.nickname,u.head_img,a.title").
		Joins("left join txy_user u on u.id = c.ouid").
		Joins("left join txy_article a on a.id = c.aid").
		Where(where).
		Limit(int(in.PageSize)).
		Offset(int(offset)).
		Order("id desc").
		Debug().
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	//计算当前type的总数，给分页算总页
	var total int64
	err = l.svcCtx.DB.Table("txy_comment as c").Where(where).Count(&total).Error
	if err != nil {
		return nil, err
	}

	return &web.CommentListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
		List:     string(jsonData),
	}, nil
}

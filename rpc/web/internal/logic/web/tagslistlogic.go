package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/pkg/model/mysql"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type TagsListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTagsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TagsListLogic {
	return &TagsListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TagsListLogic) TagsList(in *web.TagsListReq) (*web.TagsListResp, error) {
	where := map[string]interface{}{}
	var results []map[string]interface{}
	err := l.svcCtx.DB.
		Model(&mysql.TxyTag{}).
		Select("txy_tag.id,txy_tag.name, COUNT(at.aid) AS count").
		Joins("left join txy_article_tag as at on at.tid = txy_tag.id").
		Where(where).
		Group("txy_tag.id").
		Order("txy_tag.id desc").
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
	err = l.svcCtx.DB.Model(&mysql.TxyTag{}).Where(where).Count(&total).Error
	if err != nil {
		return nil, err
	}

	return &web.TagsListResp{
		List: string(jsonData),
	}, nil
}

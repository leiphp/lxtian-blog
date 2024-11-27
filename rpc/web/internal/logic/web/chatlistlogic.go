package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/rpc/web/model/mysql"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChatListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatListLogic {
	return &ChatListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChatListLogic) ChatList(in *web.ChatListReq) (*web.ChatListResp, error) {
	where := map[string]interface{}{}
	where["status"] = 1
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 3
	}
	offset := (in.Page - 1) * in.PageSize
	var results []map[string]interface{}
	err := l.svcCtx.DB.
		Model(&mysql.TxyChat{}).
		Select("id,title,cid,ctime,view_count,status").
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
	err = l.svcCtx.DB.Model(&mysql.TxyChat{}).Where(where).Count(&total).Error
	if err != nil {
		return nil, err
	}

	return &web.ChatListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
		List:     string(jsonData),
	}, nil
}

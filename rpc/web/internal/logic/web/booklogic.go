package weblogic

import (
	"context"
	"encoding/json"
	"lxtian-blog/common/pkg/model/mysql"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookLogic {
	return &BookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BookLogic) Book(in *web.BookReq) (*web.BookResp, error) {
	res := make(map[string]interface{})
	var book map[string]interface{}
	var chapters []map[string]interface{}
	err := l.svcCtx.DB.
		Model(&mysql.TxyBook{}).
		Select("id,title,slug").
		Where("id =?", in.Id).
		Order("id desc").
		Debug().
		Find(&book).Error
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.DB.
		Model(&mysql.TxyChapter{}).
		Select("id,title,parent_id,is_group,sort").
		Where("book_id = ?", in.Id).
		Order("id asc").
		Debug().
		Find(&chapters).Error
	if err != nil {
		return nil, err
	}
	res["book"] = book
	res["chapters"] = utils.BuildTree(chapters, 0)
	jsonData, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return &web.BookResp{
		Data: string(jsonData),
	}, nil
}

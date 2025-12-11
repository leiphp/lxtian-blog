package weblogic

import (
	"context"
	"encoding/json"
	"fmt"
	"lxtian-blog/rpc/web/internal/consts"
	"strings"

	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBookListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookListLogic {
	return &BookListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BookListLogic) BookList(in *web.BookListReq) (*web.BookListResp, error) {
	fmt.Println("in:", in)
	where := map[string]interface{}{}
	baseDB := l.svcCtx.DB.
		Table("txy_book as b").
		Joins("left join txy_column c on c.id = b.column_id")
	where["status"] = consts.BookStatusActive
	if in.Page == 0 {
		in.Page = 1
	}
	if in.PageSize == 0 {
		in.PageSize = 10
	}
	if in.Keywords != "" {
		baseDB = baseDB.Where("b.title like ?", "%"+in.Keywords+"%")
	}
	if in.Column > 0 {
		where["column_id"] = in.Column
	}
	offset := (in.Page - 1) * in.PageSize
	var results []map[string]interface{}
	err := baseDB.Select("b.id,b.title,b.cover,b.status,c.name column_name").
		Where(where).
		Limit(int(in.PageSize)).
		Offset(int(offset)).
		Order("b.id desc").
		Debug().
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	for k, book := range results {
		if !strings.HasPrefix(book["cover"].(string), "http://") && !strings.HasPrefix(book["cover"].(string), "https://") {
			results[k]["cover"] = l.svcCtx.QiniuClient.PrivateURL(book["cover"].(string), 3600)
		}
	}
	jsonData, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}
	//计算当前type的总数，给分页算总页
	var total int64
	err = baseDB.Count(&total).Error
	if err != nil {
		return nil, err
	}

	return &web.BookListResp{
		Page:     in.Page,
		PageSize: in.PageSize,
		Total:    uint32(total),
		List:     string(jsonData),
	}, nil
}

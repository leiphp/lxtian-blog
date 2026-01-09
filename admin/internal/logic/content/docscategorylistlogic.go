package content

import (
	"context"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/utils"
	"lxtian-blog/common/repository/web_repo"

	"github.com/zeromicro/go-zero/core/logx"
)

type DocsCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 专栏列表
func NewDocsCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsCategoryListLogic {
	return &DocsCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DocsCategoryListLogic) DocsCategoryList() (resp *types.DocsCategoryListResp, err error) {
	repo := web_repo.NewTxyDocsCategoriesRepository(l.svcCtx.DB)
	// 构建查询条件
	condition := make(map[string]interface{})
	result, _, err := repo.GetList(l.ctx, condition, 1, 999, "", "")
	if err != nil {
		return nil, err
	}

	// 转换为 map 便于前端消费
	list, err := utils.StructSliceToMapSliceUsingJSON(result)
	if err != nil {
		return nil, err
	}

	utils.FormatTimeFields(list, "created_at", "updated_at")
	resp = &types.DocsCategoryListResp{
		Data: list,
	}
	return
}

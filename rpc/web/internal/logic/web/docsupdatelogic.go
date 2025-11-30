package weblogic

import (
	"context"
	"errors"
	"fmt"
	"lxtian-blog/common/model"
	"lxtian-blog/common/pkg/redis"
	"lxtian-blog/rpc/web/internal/svc"
	"lxtian-blog/rpc/web/web"
	"time"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DocsUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDocsUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DocsUpdateLogic {
	return &DocsUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DocsUpdateLogic) DocsUpdate(in *web.DocsUpdateReq) (*web.DocsUpdateResp, error) {
	// 检查文档是否存在
	docID := int32(in.Id)
	var doc model.TxyDoc
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ?", docID).
		First(&doc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("文档不存在: id=%d", docID)
		}
		logc.Errorf(l.ctx, "查询文档失败: %s", err)
		return nil, fmt.Errorf("查询文档失败: %w", err)
	}

	// 更新文档内容
	updateData := model.TxyDoc{
		Content: &in.Content,
	}
	now := time.Now()
	updateData.UpdatedAt = &now

	// 执行更新
	err = l.svcCtx.DB.WithContext(l.ctx).
		Model(&model.TxyDoc{}).
		Where("id = ?", docID).
		Updates(updateData).Error
	if err != nil {
		logc.Errorf(l.ctx, "更新文档失败: %s", err)
		return nil, fmt.Errorf("更新文档失败: %w", err)
	}

	// 清除缓存
	cacheKey := redis.ReturnRedisKey(redis.ApiWebStringDocDetail, docID)
	if _, err = l.svcCtx.Rds.Del(cacheKey); err != nil {
		logc.Errorf(l.ctx, "清除文档缓存失败: %s", err)
		// 缓存清除失败不影响更新结果
	}

	logx.Infof("文档 %d 更新成功", docID)

	return &web.DocsUpdateResp{
		Status: "success",
		Id:     in.Id,
	}, nil
}

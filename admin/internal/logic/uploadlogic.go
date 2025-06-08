package logic

import (
	"context"
	"fmt"
	"lxtian-blog/admin/internal/svc"
	"lxtian-blog/admin/internal/types"
	"lxtian-blog/common/pkg/utils"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 图片上传
func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadLogic) Upload(r *http.Request) (resp *types.UploadResp, err error) {
	resp = new(types.UploadResp)
	// 获取上传文件
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 自定义保存逻辑，例如保存到七牛或本地
	filename := utils.GenerateFilename(header.Filename)
	// 上传文件
	url, err := l.svcCtx.QiniuClient.UploadFile(file, fmt.Sprintf("blog/cover/%s", filename))
	if err != nil {
		return nil, err
	}
	resp.Url = url
	return
}

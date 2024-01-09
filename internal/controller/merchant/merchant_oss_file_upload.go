package merchant

import (
	"context"
	ossService "go-oversea-pay/internal/logic/oss"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"go-oversea-pay/api/merchant/oss"
)

func (c *ControllerOss) FileUpload(ctx context.Context, req *oss.FileUploadReq) (res *oss.FileUploadRes, err error) {
	if req.File == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "请选择需要上传的文件")
	}
	result, err := ossService.Upload(ctx, ossService.FileUploadInput{
		File:       req.File,
		RandomName: true,
	})
	if err != nil {
		return nil, err
	}
	res = &oss.FileUploadRes{
		Url: result.Url,
	}
	return res, nil
}

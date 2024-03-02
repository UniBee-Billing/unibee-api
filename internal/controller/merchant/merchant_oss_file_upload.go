package merchant

import (
	"context"
	ossService "unibee/internal/logic/oss"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/oss"
)

func (c *ControllerOss) FileUpload(ctx context.Context, req *oss.FileUploadReq) (res *oss.FileUploadRes, err error) {
	if req.File == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "Please Specify The File")
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

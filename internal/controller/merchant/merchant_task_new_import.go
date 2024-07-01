package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/batch"
	ossService "unibee/internal/logic/oss"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/task"
)

func (c *ControllerTask) NewImport(ctx context.Context, req *task.NewImportReq) (res *task.NewImportRes, err error) {
	utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "No Permission")
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
	err = batch.NewBatchImportTask(ctx, &batch.MerchantBatchImportTaskInternalRequest{
		MerchantId:    _interface.GetMerchantId(ctx),
		MemberId:      _interface.Context().Get(ctx).MerchantMember.Id,
		Task:          req.Task,
		UploadFileUrl: result.Url,
	})
	if err != nil {
		return nil, err
	}
	return &task.NewImportRes{}, nil
}

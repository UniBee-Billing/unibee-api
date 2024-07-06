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
	taskImpl := batch.GetImportTaskImpl(req.Task)
	utility.Assert(taskImpl != nil, "Task not found")
	result, err := ossService.Upload(ctx, ossService.FileUploadInput{
		File:       req.File,
		RandomName: true,
	})
	if err != nil {
		return nil, err
	}
	////version control
	//{
	//	importFile := utility.DownloadFile(result.Url)
	//	if len(importFile) == 0 {
	//		return nil, gerror.Newf("download url failed:%s", result.Url)
	//	}
	//	reader, err := excelize.OpenFile(importFile)
	//	defer func() {
	//		if err := reader.Close(); err != nil {
	//			fmt.Println(err)
	//		}
	//	}()
	//	if err != nil {
	//		return nil, err
	//	}
	//	cell, err := reader.GetCellValue(taskImpl.TaskName(), "BB1")
	//	if err != nil {
	//		return nil, err
	//	}
	//	utility.Assert(cell == taskImpl.TemplateVersion(), "Template file is deprecated, please re-download")
	//}
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

package batch

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xuri/excelize/v2"
	"strconv"
	"unibee/internal/cmd/config"
	"unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/batch/task/invoice"
	"unibee/internal/logic/oss"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

var taskMap = map[string]_interface.BatchTask{
	"InvoiceExport": &invoice.TaskInvoice{},
}

func GetTask(task string) _interface.BatchTask {
	return taskMap[task]
}

type MerchantBatchTaskInternalRequest struct {
	MerchantId uint64            `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	MemberId   uint64            `json:"memberId" dc:"MemberId" `
	Task       string            `json:"task" dc:"Task"`
	Payload    map[string]string `json:"payload" dc:"Payload"`
}

func NewBatchDownloadTask(superCtx context.Context, req *MerchantBatchTaskInternalRequest) error {
	if len(config.GetConfigInstance().MinioConfig.Endpoint) == 0 ||
		len(config.GetConfigInstance().MinioConfig.BucketName) == 0 ||
		len(config.GetConfigInstance().MinioConfig.AccessKey) == 0 ||
		len(config.GetConfigInstance().MinioConfig.SecretKey) == 0 {
		g.Log().Errorf(superCtx, "UploadInvoicePdf error:Oss service not setup")
		utility.Assert(true, "File service need setup")
	}
	utility.Assert(req.MerchantId > 0, "Invalid Merchant")
	utility.Assert(req.MemberId > 0, "Invalid Member")
	utility.Assert(len(req.Task) > 0, "Invalid Task")
	task := GetTask(req.Task)
	utility.Assert(task != nil, "Task not found")
	one := &entity.MerchantBatchTask{
		MerchantId:   req.MerchantId,
		MemberId:     req.MemberId,
		ModuleName:   "",
		TaskName:     task.TableName(),
		SourceFrom:   "",
		Payload:      utility.MarshalToJsonString(req.Payload),
		Status:       0,
		StartTime:    0,
		FinishTime:   0,
		TaskCost:     0,
		FailReason:   "",
		GmtCreate:    nil,
		TaskType:     0,
		SuccessCount: 0,
		CreateTime:   gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantBatchTask.Ctx(superCtx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`BatchDownloadTask record insert failure %s`, err.Error())
		return err
	}
	id, _ := result.LastInsertId()
	one.Id = int64(uint(id))
	utility.Assert(one.Id > 0, "BatchDownloadTask record insert failure")
	StartRunTaskBackground(one, task)
	return nil
}

func StartRunTaskBackground(one *entity.MerchantBatchTask, task _interface.BatchTask) {
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				log.PrintPanic(ctx, err)
				return
			}
		}()
		file := excelize.NewFile()
		var startTime = gtime.Now().Timestamp()
		_, err = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
			dao.MerchantBatchTask.Columns().Status:       1,
			dao.MerchantBatchTask.Columns().StartTime:    startTime,
			dao.MerchantBatchTask.Columns().FinishTime:   0,
			dao.MerchantBatchTask.Columns().TaskCost:     0,
			dao.MerchantBatchTask.Columns().SuccessCount: 0,
			dao.MerchantBatchTask.Columns().FailReason:   "",
			dao.MerchantBatchTask.Columns().GmtModify:    gtime.Now(),
		}).Where(dao.MerchantBatchTask.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			FailureTask(ctx, one.Id, err)
			return
		}

		//Set Header
		err = file.SetSheetName("Sheet1", task.TableName())
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			FailureTask(ctx, one.Id, err)
			return
		}
		//Create Stream Writer
		writer, err := file.NewStreamWriter(task.TableName())
		//Update Width Height
		err = writer.SetColWidth(1, 15, 12)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			FailureTask(ctx, one.Id, err)
			return
		}
		//Set Header
		err = writer.SetRow("A1", task.Header())
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			FailureTask(ctx, one.Id, err)
			return
		}
		var page = 0
		var count = 100
		for {
			data, pageDataErr := task.PageData(ctx, page, count, one)
			if pageDataErr != nil {
				FailureTask(ctx, one.Id, pageDataErr)
				return
			}
			if data == nil {
				break
			}
			for i, one := range data {
				cell, _ := excelize.CoordinatesToCellName(1, page*count+i+1)
				_ = writer.SetRow(cell, one)
			}
			err = writer.Flush()
			if err != nil {
				g.Log().Errorf(ctx, err.Error())
				FailureTask(ctx, one.Id, err)
				return
			}
			_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
				dao.MerchantBatchTask.Columns().SuccessCount:   gdb.Raw(fmt.Sprintf("success_count + %v", len(data))),
				dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
				dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
			}).Where(dao.MerchantBatchTask.Columns().Id, one.Id).OmitNil().Update()
			if len(data) < count {
				break
			}
			page = page + 1
		}
		fileName := fmt.Sprintf("Batch_task_%v_%v_%v.xlsx", one.MerchantId, one.MerchantId, one.Id)
		err = file.SaveAs(fileName)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			FailureTask(ctx, one.Id, err)
			return
		}
		upload, err := oss.UploadLocalFile(ctx, fileName, "batch_download", fmt.Sprintf("%v_%v_%v", one.MerchantId, one.MemberId, one.Id), strconv.FormatUint(one.MemberId, 10))
		if err != nil {
			g.Log().Errorf(ctx, fmt.Sprintf("StartRunTaskBackground UploadLocalFile error:%v", err))
			return
		}
		_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
			dao.MerchantBatchTask.Columns().Status:         2,
			dao.MerchantBatchTask.Columns().UploadFileUrl:  upload.Url,
			dao.MerchantBatchTask.Columns().FinishTime:     gtime.Now().Timestamp(),
			dao.MerchantBatchTask.Columns().TaskCost:       gtime.Now().Timestamp() - startTime,
			dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
			dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
		}).Where(dao.MerchantBatchTask.Columns().Id, one.Id).OmitNil().Update()
	}()
}

func FailureTask(ctx context.Context, taskId int64, err error) {
	_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
		dao.MerchantBatchTask.Columns().Status:         3,
		dao.MerchantBatchTask.Columns().FailReason:     fmt.Sprintf("%s", err.Error()),
		dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
		dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
	}).Where(dao.MerchantBatchTask.Columns().Id, taskId).OmitNil().Update()
}

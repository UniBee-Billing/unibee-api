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
	"unibee/internal/logic/batch/export/discount"
	"unibee/internal/logic/batch/export/invoice"
	"unibee/internal/logic/batch/export/subscription"
	"unibee/internal/logic/batch/export/transaction"
	"unibee/internal/logic/batch/export/user"
	"unibee/internal/logic/oss"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

var exportTaskMap = map[string]_interface.BatchExportTask{
	"InvoiceExport":      &invoice.TaskInvoiceExport{},
	"UserExport":         &user.TaskUserExport{},
	"SubscriptionExport": &subscription.TaskSubscriptionExport{},
	"TransactionExport":  &transaction.TaskTransactionExport{},
	"DiscountExport":     &discount.TaskDiscountExport{},
	"UserDiscountExport": &discount.TaskUserDiscountExport{},
}

func GetExportTaskImpl(task string) _interface.BatchExportTask {
	return exportTaskMap[task]
}

type MerchantBatchExportTaskInternalRequest struct {
	MerchantId        uint64                 `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	MemberId          uint64                 `json:"memberId" dc:"MemberId" `
	Task              string                 `json:"task" dc:"Task"`
	Payload           map[string]interface{} `json:"payload" dc:"Payload"`
	SkipColumnIndexes []int                  `json:"skipColumnIndexes" dc:"SkipColumnIndexes, the column will be skipped in the export file if its index specified"`
}

func ExportColumnList(ctx context.Context, task string) []interface{} {
	one := GetExportTaskImpl(task)
	utility.Assert(one != nil, "Task not found")
	return RefactorHeaders(one.Header(), nil)
}

func NewBatchExportTask(superCtx context.Context, req *MerchantBatchExportTaskInternalRequest) error {
	if len(config.GetConfigInstance().MinioConfig.Endpoint) == 0 ||
		len(config.GetConfigInstance().MinioConfig.BucketName) == 0 ||
		len(config.GetConfigInstance().MinioConfig.AccessKey) == 0 ||
		len(config.GetConfigInstance().MinioConfig.SecretKey) == 0 {
		g.Log().Errorf(superCtx, "NewBatchExportTask error:file service not setup")
		utility.Assert(true, "File service need setup")
	}
	utility.Assert(req.MerchantId > 0, "Invalid Merchant")
	utility.Assert(req.MemberId > 0, "Invalid Member")
	utility.Assert(len(req.Task) > 0, "Invalid Task")
	task := GetExportTaskImpl(req.Task)
	utility.Assert(task != nil, "Task not found")
	one := &entity.MerchantBatchTask{
		MerchantId:   req.MerchantId,
		MemberId:     req.MemberId,
		ModuleName:   "",
		TaskName:     task.TaskName(),
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
		err = gerror.Newf(`BatchExportTask record insert failure %s`, err.Error())
		return err
	}
	id, _ := result.LastInsertId()
	one.Id = int64(uint(id))
	utility.Assert(one.Id > 0, "BatchExportTask record insert failure")
	startRunExportTaskBackground(one, task, req.SkipColumnIndexes)
	return nil
}

func startRunExportTaskBackground(task *entity.MerchantBatchTask, taskImpl _interface.BatchExportTask, skipColumnIndexes []int) {
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
				failureTask(ctx, task.Id, err)
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
		}).Where(dao.MerchantBatchTask.Columns().Id, task.Id).OmitNil().Update()
		if err != nil {
			failureTask(ctx, task.Id, err)
			return
		}

		//Set Header
		err = file.SetSheetName("Sheet1", taskImpl.TaskName())
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		//Create Stream Writer
		writer, err := file.NewStreamWriter(taskImpl.TaskName())
		//Update Width Height
		err = writer.SetColWidth(1, 15, 12)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		headerStyleID, err := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}

		err = writer.SetRow("A1", RefactorHeaders(taskImpl.Header(), skipColumnIndexes), excelize.RowOpts{StyleID: headerStyleID})
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		var page = 0
		var count = 100
		for {
			list, pageDataErr := taskImpl.PageData(ctx, page, count, task)
			if pageDataErr != nil {
				failureTask(ctx, task.Id, pageDataErr)
				return
			}
			if list == nil {
				break
			}
			for i, one := range list {
				if one == nil {
					continue
				}
				cell, _ := excelize.CoordinatesToCellName(1, page*count+i+2)
				_ = writer.SetRow(cell, RefactorData(one, "", skipColumnIndexes))
			}
			_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
				dao.MerchantBatchTask.Columns().SuccessCount:   gdb.Raw(fmt.Sprintf("success_count + %v", len(list))),
				dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
				dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
			}).Where(dao.MerchantBatchTask.Columns().Id, task.Id).OmitNil().Update()
			if len(list) < count {
				break
			}
			page = page + 1
		}
		err = writer.Flush()
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		fileName := fmt.Sprintf("Batch_export_task_%v_%v_%v.xlsx", task.MerchantId, task.MemberId, task.Id)
		err = file.SaveAs(fileName)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		upload, err := oss.UploadLocalFile(ctx, fileName, "batch_export", fileName, strconv.FormatUint(task.MemberId, 10))
		if err != nil {
			g.Log().Errorf(ctx, fmt.Sprintf("startRunExportTaskBackground UploadLocalFile error:%v", err))
			failureTask(ctx, task.Id, err)
			return
		}
		_, err = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
			dao.MerchantBatchTask.Columns().Status:         2,
			dao.MerchantBatchTask.Columns().DownloadUrl:    upload.Url,
			dao.MerchantBatchTask.Columns().FinishTime:     gtime.Now().Timestamp(),
			dao.MerchantBatchTask.Columns().TaskCost:       gtime.Now().Timestamp() - startTime,
			dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
			dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
		}).Where(dao.MerchantBatchTask.Columns().Id, task.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(ctx, fmt.Sprintf("startRunExportTaskBackground Update MerchantBatchTask error:%v", err))
			failureTask(ctx, task.Id, err)
			return
		}
	}()
}

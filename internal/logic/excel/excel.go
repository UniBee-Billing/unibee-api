package excel

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xuri/excelize/v2"
	"unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/excel/task/invoice"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

var taskMap = map[string]_interface.BatchTask{
	"InvoiceExport": &invoice.TaskInvoice{},
}

func GetTaskForName(taskName string) _interface.BatchTask {
	return taskMap[taskName]
}

type MerchantBatchTaskInternalRequest struct {
	MerchantId uint64            `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	MemberId   uint64            `json:"memberId" dc:"MemberId" `
	TaskName   string            `json:"taskName" dc:"TaskName"`
	ModuleName string            `json:"moduleName" dc:"ModuleName"`
	SourceFrom string            `json:"sourceFrom" dc:"SourceFrom"`
	Payload    map[string]string `json:"payload" dc:"Payload"`
}

func BatchDownloadTask(superCtx context.Context, req MerchantBatchTaskInternalRequest) error {
	utility.Assert(req.MerchantId > 0, "Invalid Merchant")
	utility.Assert(req.MemberId > 0, "Invalid Member")
	utility.Assert(len(req.TaskName) > 0, "Invalid TaskName")
	task := GetTaskForName(req.TaskName)
	utility.Assert(task != nil, "Invalid TaskName")
	one := &entity.MerchantBatchTask{
		MerchantId:   req.MerchantId,
		MemberId:     req.MemberId,
		ModuleName:   req.ModuleName,
		TaskName:     req.TaskName,
		SourceFrom:   req.SourceFrom,
		Payload:      utility.MarshalToJsonString(req.Payload),
		Status:       0,
		StartTime:    0,
		FinishTime:   0,
		TaskCost:     0,
		FailReason:   "",
		GmtCreate:    nil,
		TaskType:     0,
		SuccessCount: 0,
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
		err = file.SetSheetName("Sheet1", task.TableName(one))
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			FailureTask(ctx, one.Id, err)
			return
		}
		//Create Stream Writer
		writer, err := file.NewStreamWriter(task.TableName(one))
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
				dao.MerchantBatchTask.Columns().SuccessCount: gdb.Raw(fmt.Sprintf("success_count + %v", len(data))),
				dao.MerchantBatchTask.Columns().GmtModify:    gtime.Now(),
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
		// todo mark upload File to Url
		_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
			dao.MerchantBatchTask.Columns().Status:     2,
			dao.MerchantBatchTask.Columns().FinishTime: gtime.Now().Timestamp(),
			dao.MerchantBatchTask.Columns().TaskCost:   gtime.Now().Timestamp() - startTime,
			dao.MerchantBatchTask.Columns().GmtModify:  gtime.Now(),
		}).Where(dao.MerchantBatchTask.Columns().Id, one.Id).OmitNil().Update()
	}()
}

func FailureTask(ctx context.Context, taskId int64, err error) {
	_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
		dao.MerchantBatchTask.Columns().Status:     3,
		dao.MerchantBatchTask.Columns().FailReason: fmt.Sprintf("%s", err.Error()),
		dao.MerchantBatchTask.Columns().GmtModify:  gtime.Now(),
	}).Where(dao.MerchantBatchTask.Columns().Id, taskId).OmitNil().Update()
}
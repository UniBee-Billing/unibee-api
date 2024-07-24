package batch

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/oversea_pay"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/batch/_import/subscription"
	user2 "unibee/internal/logic/batch/_import/user"
	"unibee/internal/logic/oss"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

var importTaskMap = map[string]_interface.BatchImportTask{
	"UserImport":                user2.TaskUserImport{},
	"ActiveSubscriptionImport":  subscription.TaskActiveSubscriptionImport{},
	"HistorySubscriptionImport": subscription.TaskHistorySubscriptionImport{},
}

func GetImportTaskImpl(task string) _interface.BatchImportTask {
	return importTaskMap[task]
}

type MerchantBatchImportTaskInternalRequest struct {
	MerchantId    uint64 `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	MemberId      uint64 `json:"memberId" dc:"MemberId" `
	Task          string `json:"task" dc:"Task"`
	UploadFileUrl string `json:"uploadFileUrl" dc:"UploadFileUrl"`
}

func NewBatchImportTask(superCtx context.Context, req *MerchantBatchImportTaskInternalRequest) error {
	if len(config.GetConfigInstance().MinioConfig.Endpoint) == 0 ||
		len(config.GetConfigInstance().MinioConfig.BucketName) == 0 ||
		len(config.GetConfigInstance().MinioConfig.AccessKey) == 0 ||
		len(config.GetConfigInstance().MinioConfig.SecretKey) == 0 {
		g.Log().Errorf(superCtx, "NewBatchImportTask error:file service not setup")
		utility.Assert(true, "File service need setup")
	}
	utility.Assert(req.MerchantId > 0, "Invalid Merchant")
	utility.Assert(req.MemberId > 0, "Invalid Member")
	utility.Assert(len(req.Task) > 0, "Invalid Task")
	task := GetImportTaskImpl(req.Task)
	utility.Assert(task != nil, "Task not found")
	one := &entity.MerchantBatchTask{
		MerchantId:    req.MerchantId,
		MemberId:      req.MemberId,
		ModuleName:    "",
		TaskName:      task.TaskName(),
		SourceFrom:    "",
		UploadFileUrl: req.UploadFileUrl,
		Payload:       "",
		Status:        0,
		StartTime:     0,
		FinishTime:    0,
		TaskCost:      0,
		FailReason:    "",
		GmtCreate:     nil,
		TaskType:      1,
		SuccessCount:  0,
		CreateTime:    gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantBatchTask.Ctx(superCtx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`BatchImportTask record insert failure %s`, err.Error())
		return err
	}
	id, _ := result.LastInsertId()
	one.Id = int64(uint(id))
	utility.Assert(one.Id > 0, "BatchImportTask record insert failure")
	startRunImportTaskBackground(one, task)
	return nil
}

func startRunImportTaskBackground(task *entity.MerchantBatchTask, taskImpl _interface.BatchImportTask) {
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
		err = file.SetSheetName("Sheet1", GeneralExportImportSheetName)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		//Create Stream Writer
		writer, err := file.NewStreamWriter(GeneralExportImportSheetName)
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
		//Set Header
		resultHeader := RefactorHeaders(taskImpl.TemplateHeader(), nil, true)
		resultHeader = append(resultHeader, "ImportResult")
		err = writer.SetRow("A1", resultHeader, excelize.RowOpts{StyleID: headerStyleID})
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		// Start Read Import File
		importFile := utility.DownloadFile(task.UploadFileUrl)
		if len(importFile) == 0 {
			err = gerror.Newf("download url failed:%s", task.UploadFileUrl)
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		reader, err := excelize.OpenFile(importFile)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		readerRows, err := reader.GetRows(GeneralExportImportSheetName)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		var headers = make(map[int]string)
		var count = 0
		for i, row := range readerRows {
			count = i
			if i == 0 {
				for j, colCell := range row {
					headers[j] = strings.Replace(colCell, " ", "", -1)
				}
			} else {
				target := make(map[string]string)
				for j, colCell := range row {
					target[headers[j]] = strings.TrimSpace(colCell)
				}
				if target == nil {
					continue
				}
				cell, _ := excelize.CoordinatesToCellName(1, i+1)
				//result, importResult := taskImpl.ImportRow(ctx, task, target)
				result, importResult := ProxyImportRow(ctx, taskImpl, task, target)
				var resultMessage = "success"
				if importResult != nil {
					resultMessage = fmt.Sprintf("%s", importResult.Error())
				}
				_ = writer.SetRow(cell, append(RefactorData(result, "", nil), resultMessage))
				if count%10 == 0 {
					_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
						dao.MerchantBatchTask.Columns().SuccessCount:   fmt.Sprintf("%v", count),
						dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
						dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
					}).Where(dao.MerchantBatchTask.Columns().Id, task.Id).OmitNil().Update()
				}
			}
		}

		_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
			dao.MerchantBatchTask.Columns().SuccessCount:   fmt.Sprintf("%v", count),
			dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
			dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
		}).Where(dao.MerchantBatchTask.Columns().Id, task.Id).OmitNil().Update()

		err = writer.Flush()
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		fileName := fmt.Sprintf("Batch_import_task_%v_%v_%v.xlsx", task.MerchantId, task.MemberId, task.Id)
		err = file.SaveAs(fileName)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		upload, err := oss.UploadLocalFile(ctx, fileName, "batch_import", fileName, strconv.FormatUint(task.MemberId, 10))
		if err != nil {
			g.Log().Errorf(ctx, fmt.Sprintf("startRunImportTaskBackground UploadLocalFile error:%v", err))
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
			g.Log().Errorf(ctx, fmt.Sprintf("startRunImportTaskBackground Update MerchantBatchTask error:%v", err))
			failureTask(ctx, task.Id, err)
			return
		}
	}()
}

func ProxyImportRow(ctx context.Context, taskImpl _interface.BatchImportTask, task *entity.MerchantBatchTask, target map[string]string) (data interface{}, err error) {
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
	return taskImpl.ImportRow(ctx, task, target)
}

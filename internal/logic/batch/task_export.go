package batch

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xuri/excelize/v2"
	"os"
	"reflect"
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
}

func GetExportTaskImpl(task string) _interface.BatchExportTask {
	return exportTaskMap[task]
}

type MerchantBatchExportTaskInternalRequest struct {
	MerchantId    uint64                 `json:"merchantId" dc:"MerchantId" v:"MerchantId"`
	MemberId      uint64                 `json:"memberId" dc:"MemberId" `
	Task          string                 `json:"task" dc:"Task"`
	Payload       map[string]interface{} `json:"payload" dc:"Payload"`
	ExportColumns []string               `json:"exportColumns" dc:"ExportColumns, the export file column list"`
	Format        string                 `json:"format" dc:"The format of export file, xlsx|csv, will be xlsx if not specified"`
}

func refactorHeaderCommentMap(obj interface{}) map[string]string {
	out := make(map[string]string, 0)
	if obj == nil {
		return out
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	utility.Assert(v.Kind() == reflect.Struct, fmt.Sprintf("ReflectTemplateStructToMap only accepts struct or struct pointer; got %T", v))

	t := v.Type()
	// range properties
	// get Tag named "json" as key
	//var allKeys = make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if key := fi.Tag.Get("json"); key != "" {
			out[key] = fi.Tag.Get("comment")
		}
	}

	return out
}

func ExportColumnList(ctx context.Context, task string) ([]interface{}, map[string]string) {
	one := GetExportTaskImpl(task)
	utility.Assert(one != nil, "Task not found")
	return RefactorHeaders(one.Header(), nil), refactorHeaderCommentMap(one.Header())
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
	if len(req.Format) > 0 {
		utility.Assert(req.Format == "xlsx" || req.Format == "csv", "format should be one of xlsx|csv")
	}
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
	startRunExportTaskBackground(one, task, req.ExportColumns, req.Format)
	return nil
}

func startRunExportTaskBackground(task *entity.MerchantBatchTask, taskImpl _interface.BatchExportTask, exportColumns []string, format string) {
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

		err = writer.SetRow("A1", RefactorHeaders(taskImpl.Header(), exportColumns), excelize.RowOpts{StyleID: headerStyleID})
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		var page = 0
		var count = 100

		//task.Payload = SupportUpperCasePayload(ctx, task.Payload)
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
				_ = writer.SetRow(cell, RefactorData(one, "", exportColumns))
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
		// addComments
		for _, comment := range RefactorHeaderComments(taskImpl.Header(), exportColumns) {
			err = file.AddComment(GeneralExportImportSheetName, comment)
		}

		fileName := fmt.Sprintf("Batch_export_task_%v_%v_%v.xlsx", task.MerchantId, task.MemberId, task.Id)
		err = file.SaveAs(fileName)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		if format == "csv" {
			fileName, err = convertXlsxFileToCSV(fileName, GeneralExportImportSheetName, fmt.Sprintf("Batch_export_task_%v_%v_%v.csv", task.MerchantId, task.MemberId, task.Id))
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

func SupportUpperCasePayload(ctx context.Context, payload string) string {
	var target map[string]interface{}
	err := utility.UnmarshalFromJsonString(payload, &target)
	if err != nil {
		g.Log().Errorf(ctx, "Download PageData error:%s", err.Error())
		return payload
	}
	for key, value := range target {
		if utility.IsStartUpper(key) {
			lowerKey := utility.ToFirstCharLowerCase(key)
			if _, ok := target[lowerKey]; !ok {
				target[lowerKey] = value
			}
		}
	}
	return utility.MarshalToJsonString(target)
}

func convertXlsxFileToCSV(xlsxFileLocalName string, sheetName string, csvFileName string) (string, error) {
	xlsxFile := xlsxFileLocalName
	f, err := excelize.OpenFile(xlsxFile)
	if err != nil {
		return "", err
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", err
	}

	csvFile := csvFileName
	csvFileHandle, err := os.Create(csvFile)
	if err != nil {
		return "", err
	}
	defer func(csvFileHandle *os.File) {
		err = csvFileHandle.Close()
		if err != nil {

		}
	}(csvFileHandle)

	writer := csv.NewWriter(csvFileHandle)
	defer writer.Flush()

	for _, row := range rows {
		if err = writer.Write(row); err != nil {
			return "", err
		}
	}
	return csvFileName, nil
}

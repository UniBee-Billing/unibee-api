package batch

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"unibee/internal/cmd/config"
	"unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/oss"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xuri/excelize/v2"
)

func NewBatchExportTaskV2(superCtx context.Context, req *MerchantBatchExportTaskInternalRequest) error {
	if len(config.GetConfigInstance().MinioConfig.Endpoint) == 0 ||
		len(config.GetConfigInstance().MinioConfig.BucketName) == 0 ||
		len(config.GetConfigInstance().MinioConfig.AccessKey) == 0 ||
		len(config.GetConfigInstance().MinioConfig.SecretKey) == 0 {
		g.Log().Errorf(superCtx, "NewBatchExportTaskV2 error:file service not setup")
		utility.Assert(true, "File service need setup")
	}
	utility.Assert(req.MerchantId > 0, "Invalid Merchant")
	utility.Assert(req.MemberId > 0, "Invalid Member")
	utility.Assert(len(req.Task) > 0, "Invalid Task")
	task := GetExportTaskImpl(req.Task)
	utility.Assert(task != nil, "Task not found")
	if len(req.Format) > 0 {
		utility.Assert(req.Format == "xlsx" || req.Format == "csv", "format should be one of xlsx|csv")
	} else {
		// Default to CSV in V2
		req.Format = "csv"
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
		Format:       req.Format,
	}
	result, err := dao.MerchantBatchTask.Ctx(superCtx).Data(one).OmitNil().Insert(one)
	if err != nil {
		err = gerror.Newf(`BatchExportTask V2 record insert failure %s`, err.Error())
		return err
	}
	id, _ := result.LastInsertId()
	one.Id = int64(uint(id))
	utility.Assert(one.Id > 0, "BatchExportTask V2 record insert failure")
	startRunExportTaskBackgroundV2(one, task, req.ExportColumns, req.Format)
	return nil
}

func startRunExportTaskBackgroundV2(task *entity.MerchantBatchTask, taskImpl _interface.BatchExportTask, exportColumns []string, format string) {
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

		// Prepare header (human-readable), store rows as chunks in DB first
		headerInterfaces := RefactorHeaders(taskImpl.Header(), exportColumns, true)
		// Clean existing chunks for this task if any
		_, _ = dao.MerchantBatchTaskExportChunks.Ctx(ctx).
			Where(dao.MerchantBatchTaskExportChunks.Columns().TaskId, uint64(task.Id)).
			Delete()

		// Write data page by page into merchant_batch_task_export_chunks
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
			// Build CSV chunk for this page
			var buf bytes.Buffer
			pageWriter := csv.NewWriter(&buf)
			for _, one := range list {
				if one == nil {
					continue
				}
				rowInterfaces := RefactorData(one, "", exportColumns)
				row := make([]string, 0, len(rowInterfaces))
				for _, v := range rowInterfaces {
					row = append(row, fmt.Sprint(v))
				}
				if err = pageWriter.Write(row); err != nil {
					failureTask(ctx, task.Id, err)
					return
				}
			}
			pageWriter.Flush()
			if err = pageWriter.Error(); err != nil {
				failureTask(ctx, task.Id, err)
				return
			}
			// Persist chunk
			_, err = dao.MerchantBatchTaskExportChunks.Ctx(ctx).Data(g.Map{
				dao.MerchantBatchTaskExportChunks.Columns().MerchantId: task.MerchantId,
				dao.MerchantBatchTaskExportChunks.Columns().TaskId:     uint64(task.Id),
				dao.MerchantBatchTaskExportChunks.Columns().Page:       page,
				dao.MerchantBatchTaskExportChunks.Columns().Content:    buf.String(),
				dao.MerchantBatchTaskExportChunks.Columns().CreateTime: gtime.Now().Timestamp(),
			}).Insert()
			if err != nil {
				failureTask(ctx, task.Id, err)
				return
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

		// Assemble final CSV from chunks
		timestamp := gtime.Now().Format("2006-01-02_15-04")
		csvFileName := fmt.Sprintf("Batch_export_%v_%s.csv", task.Id, timestamp)
		csvFileHandle, err := os.Create(csvFileName)
		if err != nil {
			g.Log().Errorf(ctx, err.Error())
			failureTask(ctx, task.Id, err)
			return
		}
		// Write header once
		header := make([]string, 0, len(headerInterfaces))
		for _, h := range headerInterfaces {
			header = append(header, fmt.Sprint(h))
		}
		headerWriter := csv.NewWriter(csvFileHandle)
		if err = headerWriter.Write(header); err != nil {
			_ = csvFileHandle.Close()
			failureTask(ctx, task.Id, err)
			return
		}
		headerWriter.Flush()
		if err = headerWriter.Error(); err != nil {
			_ = csvFileHandle.Close()
			failureTask(ctx, task.Id, err)
			return
		}
		// Stream chunks in order
		var chunks []entity.MerchantBatchTaskExportChunks
		err = dao.MerchantBatchTaskExportChunks.Ctx(ctx).
			Where(dao.MerchantBatchTaskExportChunks.Columns().TaskId, uint64(task.Id)).
			OrderAsc(dao.MerchantBatchTaskExportChunks.Columns().Page).
			Scan(&chunks)
		if err != nil {
			_ = csvFileHandle.Close()
			failureTask(ctx, task.Id, err)
			return
		}
		for _, ch := range chunks {
			if len(ch.Content) == 0 {
				continue
			}
			if _, err = csvFileHandle.WriteString(ch.Content); err != nil {
				_ = csvFileHandle.Close()
				failureTask(ctx, task.Id, err)
				return
			}
		}
		if err = csvFileHandle.Close(); err != nil {
			failureTask(ctx, task.Id, err)
			return
		}

		finalFileName := csvFileName
		if format == "xlsx" {
			headerComments := RefactorHeaderComments(taskImpl.Header(), exportColumns)
			xlsxFileName := fmt.Sprintf("Batch_export_%v_%s.xlsx", task.Id, timestamp)
			if finalFileName, err = convertCSVToXlsx(csvFileName, xlsxFileName, headerComments); err != nil {
				g.Log().Errorf(ctx, err.Error())
				failureTask(ctx, task.Id, err)
				return
			}
		}

		upload, err := oss.UploadLocalFile(ctx, finalFileName, "batch_export", finalFileName, strconv.FormatUint(task.MemberId, 10))
		if err != nil {
			g.Log().Errorf(ctx, fmt.Sprintf("startRunExportTaskBackgroundV2 UploadLocalFile error:%v", err))
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
			g.Log().Errorf(ctx, fmt.Sprintf("startRunExportTaskBackgroundV2 Update MerchantBatchTask error:%v", err))
			failureTask(ctx, task.Id, err)
			return
		}
		// Cleanup cached chunks after file uploaded and task marked success
		if _, delErr := dao.MerchantBatchTaskExportChunks.Ctx(ctx).
			Where(dao.MerchantBatchTaskExportChunks.Columns().TaskId, uint64(task.Id)).
			Delete(); delErr != nil {
			g.Log().Warningf(ctx, "cleanup MerchantBatchTaskExportChunks failed, taskId=%d, err=%v", task.Id, delErr)
		}
	}()
}

func convertCSVToXlsx(csvFileName string, xlsxFileName string, headerComments []excelize.Comment) (string, error) {
	// Read CSV and stream-write to XLSX to avoid huge memory usage
	f, err := os.Open(csvFileName)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	reader := csv.NewReader(f)
	xl := excelize.NewFile()
	if err := xl.SetSheetName("Sheet1", GeneralExportImportSheetName); err != nil {
		return "", err
	}
	stream, err := xl.NewStreamWriter(GeneralExportImportSheetName)
	if err != nil {
		return "", err
	}
	// Align with v1: set column width and bold header row
	_ = stream.SetColWidth(1, 15, 12)
	headerStyleID, _ := xl.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})

	rowIndex := 1
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
		row := make([]interface{}, 0, len(record))
		for _, v := range record {
			row = append(row, v)
		}
		if rowIndex == 1 {
			if err = stream.SetRow(cell, row, excelize.RowOpts{StyleID: headerStyleID}); err != nil {
				return "", err
			}
		} else {
			if err = stream.SetRow(cell, row); err != nil {
				return "", err
			}
		}
		rowIndex++
	}
	if err := stream.Flush(); err != nil {
		return "", err
	}
	// Add header comments back to first row like v1
	for _, comment := range headerComments {
		_ = xl.AddComment(GeneralExportImportSheetName, comment)
	}
	if err := xl.SaveAs(xlsxFileName); err != nil {
		return "", err
	}
	return xlsxFileName, nil
}

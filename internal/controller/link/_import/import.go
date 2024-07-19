package _import

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
	"reflect"
	"unibee/internal/logic/batch"
	"unibee/utility"
)

func LinkImportTemplateEntry(r *ghttp.Request) {
	g.Log().Infof(r.Context(), "LinkImportTemplateEntry:%v", r.Method)
	r.Response.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	r.Response.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT,DELETE,OPTIONS,PATCH")
	r.Response.Header().Add("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		return
	}

	task := r.Get("task").String()
	if len(task) <= 0 {
		r.Response.Writeln("Task invalid")
		return
	}

	taskImpl := batch.GetImportTaskImpl(utility.Case2Camel(task))
	if taskImpl == nil {
		r.Response.Writeln("Task not found")
		return
	}
	file := excelize.NewFile()
	err := file.SetSheetName("Sheet1", batch.GeneralExportImportSheetName)
	if err != nil {
		g.Log().Errorf(r.Context(), err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}
	//Create Stream Writer
	writer, err := file.NewStreamWriter(batch.GeneralExportImportSheetName)
	//Update Width Height
	err = writer.SetColWidth(1, 15, 12)
	if err != nil {
		g.Log().Errorf(r.Context(), err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}
	//Set Header
	headerStyleID, err := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})
	if err != nil {
		g.Log().Errorf(r.Context(), err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}
	err = writer.SetRow("A1", batch.RefactorHeaders(taskImpl.TemplateHeader(), nil), excelize.RowOpts{StyleID: headerStyleID})
	if err != nil {
		g.Log().Errorf(r.Context(), err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}
	cell, _ := excelize.CoordinatesToCellName(1, 2)
	_ = writer.SetRow(cell, batch.RefactorData(taskImpl.TemplateHeader(), "", nil))

	err = writer.Flush()
	if err != nil {
		g.Log().Errorf(r.Context(), err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}
	fileName := fmt.Sprintf("ImportTemplate_%v.xlsx", task)

	// addComments
	for _, comment := range refactorHeaderComments(taskImpl.TemplateHeader(), nil) {
		err = file.AddComment(batch.GeneralExportImportSheetName, comment)
	}

	err = file.SaveAs(fileName)
	if err != nil {
		g.Log().Errorf(r.Context(), err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}

	if len(fileName) == 0 {
		g.Log().Errorf(r.Context(), "LinkEntry pdfFile download or generate error")
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}

	r.Response.Header().Add("Content-type", "application/octet-stream")
	r.Response.Header().Add("content-disposition", "attachment; filename=\""+fileName+"\"")
	downloadFile, err := os.Open(fileName)
	if err != nil {
		g.Log().Errorf(r.Context(), "LinkEntry error:%s", err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			g.Log().Errorf(r.Context(), "LinkEntry error:%s", err.Error())
		}
	}(downloadFile)

	_, err = io.Copy(r.Response.ResponseWriter, downloadFile)
	if err != nil {
		g.Log().Errorf(r.Context(), "LinkEntry error:%s", err.Error())
		r.Response.WriteHeader(http.StatusBadRequest)
		r.Response.Writeln("Bad request")
	}
}

func refactorHeaderComments(obj interface{}, skipColumnIndexes []int) []excelize.Comment {
	out := make([]excelize.Comment, 0)
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
	for i := 0; i < v.NumField(); i++ {
		if utility.IsIntInArray(skipColumnIndexes, i) {
			continue
		}
		fi := t.Field(i)
		if comment := fi.Tag.Get("comment"); comment != "" {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			out = append(out, excelize.Comment{
				Cell:   cell,
				Author: "UniBee",
				Paragraph: []excelize.RichTextRun{
					{Text: comment, Font: &excelize.Font{Bold: true}},
				},
			})
		}
	}

	return out
}

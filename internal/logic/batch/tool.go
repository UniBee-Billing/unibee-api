package batch

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/xuri/excelize/v2"
	"reflect"
	"time"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/utility"
)

const GeneralExportImportSheetName = "Sheet1"

func RefactorHeaders(obj interface{}, exportColumns []string) []interface{} {
	out := make([]interface{}, 0)
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
	var allKeys = make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if key := fi.Tag.Get("json"); key != "" {
			out = append(out, key)
			allKeys[key] = "1"
		}
	}

	if exportColumns != nil && len(exportColumns) > 0 {
		out = nil
		for _, key := range exportColumns {
			if _, ok := allKeys[key]; ok {
				out = append(out, key)
			}
		}
		return out
	} else {
		return out
	}
}

func RefactorData(obj interface{}, timeZone string, exportColumns []string) []interface{} {
	out := make([]interface{}, 0)
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
	var allValue = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if key := fi.Tag.Get("json"); key != "" {
			value := v.Field(i).Interface()
			if layout := fi.Tag.Get("layout"); layout != "" {
				if targetTime, ok := value.(*gtime.Time); ok {
					if len(timeZone) > 0 {
						loc, err := time.LoadLocation(timeZone)
						if err == nil {
							targetTime = targetTime.ToLocation(loc)
						}
					}
					value = targetTime.Layout(layout)
					if value == "0001-01-01 00:00:00" {
						value = ""
					}
				}
			}
			if value == nil {
				allValue[key] = ""
			} else {
				allValue[key] = value
			}
			out = append(out, value)
		}
	}

	if exportColumns != nil && len(exportColumns) > 0 {
		out = nil
		for _, key := range exportColumns {
			if value, ok := allValue[key]; ok {
				out = append(out, value)
			}
		}
		return out
	} else {
		return out
	}

}

func RefactorHeaderComments(obj interface{}, exportColumns []string) []excelize.Comment {
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
	var allKeys = make(map[string]string)
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if key := fi.Tag.Get("json"); key != "" {
			if comment := fi.Tag.Get("comment"); comment != "" {
				allKeys[key] = comment
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
	}

	if exportColumns != nil && len(exportColumns) > 0 {
		out = nil
		for i, key := range exportColumns {
			if comment, ok := allKeys[key]; ok {
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
	} else {
		return out
	}

}

func failureTask(ctx context.Context, taskId int64, err error) {
	_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
		dao.MerchantBatchTask.Columns().Status:         3,
		dao.MerchantBatchTask.Columns().FailReason:     fmt.Sprintf("%s \n \t %s", err.Error(), g.Log().GetStack()),
		dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
		dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
	}).Where(dao.MerchantBatchTask.Columns().Id, taskId).OmitNil().Update()
}

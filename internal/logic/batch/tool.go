package batch

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"reflect"
	"time"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/utility"
)

func RefactorHeaders(obj interface{}, skipColumnIndexes []int) []interface{} {
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
	for i := 0; i < v.NumField(); i++ {
		if utility.IsIntInArray(skipColumnIndexes, i) {
			continue
		}
		fi := t.Field(i)
		if key := fi.Tag.Get("json"); key != "" {
			out = append(out, key)
		}
	}

	return out
}

func RefactorData(obj interface{}, timeZone string, skipColumnIndexes []int) []interface{} {
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
	for i := 0; i < v.NumField(); i++ {
		if utility.IsIntInArray(skipColumnIndexes, i) {
			continue
		}
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
			out = append(out, value)
		}
	}
	return out
}

func failureTask(ctx context.Context, taskId int64, err error) {
	_, _ = dao.MerchantBatchTask.Ctx(ctx).Data(g.Map{
		dao.MerchantBatchTask.Columns().Status:         3,
		dao.MerchantBatchTask.Columns().FailReason:     fmt.Sprintf("%s \n \t %s", err.Error(), g.Log().GetStack()),
		dao.MerchantBatchTask.Columns().LastUpdateTime: gtime.Now().Timestamp(),
		dao.MerchantBatchTask.Columns().GmtModify:      gtime.Now(),
	}).Where(dao.MerchantBatchTask.Columns().Id, taskId).OmitNil().Update()
}

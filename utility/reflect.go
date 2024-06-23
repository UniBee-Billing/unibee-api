package utility

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"reflect"
	"time"
)

func ReflectTemplateStructToMap(in interface{}, timeZone string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ReflectTemplateStructToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// range properties
	// get Tag named "json" as key
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			target := v.Field(i).Interface()
			if layout := fi.Tag.Get("layout"); layout != "" {
				if targetTime, ok := target.(*gtime.Time); ok {
					if len(timeZone) > 0 {
						loc, err := time.LoadLocation(timeZone)
						if err == nil {
							targetTime = targetTime.ToLocation(loc)
						}
					}
					target = targetTime.Layout(layout)
				}
			}
			out[tagValue] = target
		}
	}
	return out, nil
}

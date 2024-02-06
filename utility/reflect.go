package utility

import (
	"fmt"
	"reflect"
)

func ReflectStructToMap(in interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { //
		return nil, fmt.Errorf("ReflectStructToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// range properties
	// get Tag named "json" as key
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get("json"); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

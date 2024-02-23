package bean

import (
	"fmt"
	"reflect"
)

func CopyProperties(source interface{}, dest interface{}) {
	sourceValue := reflect.ValueOf(source).Elem()
	destValue := reflect.ValueOf(dest).Elem()

	for i := 0; i < sourceValue.NumField(); i++ {
		fmt.Println("source Type:", sourceValue.Type())
		destFieldValue := destValue.FieldByName(sourceValue.Type().Field(i).Name)
		fmt.Println("dest Type:", destFieldValue.Type())
		if destFieldValue.IsValid() && destFieldValue.CanSet() {
			destFieldValue.Set(sourceValue.Field(i))
		}
	}
}

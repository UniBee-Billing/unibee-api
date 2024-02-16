package redismq

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"reflect"
)

func Serialize(target interface{}) []byte {
	jsonData, _ := gjson.Marshal(target)
	return jsonData
}

func Deserialize(body []byte, v interface{}) (err error) {
	if !isPointerType(v) {
		err = gerror.New("v should be pointer type")
		return
	}
	err = gjson.Unmarshal(body, &v) // Unmarshal todo mark 加上 &
	return
}

func isPointerType(value interface{}) bool {
	typ := reflect.TypeOf(value)
	kind := typ.Kind()

	return kind == reflect.Ptr
}

//func Deserialize(body []byte) interface{} {
//	var result interface{}
//	err := gjson.Unmarshal(body, &result)
//	if err != nil {
//		fmt.Printf("Deserialize err:%s\n", err)
//	}
//	return result
//}

// test case below

type Person struct {
	Name   string
	Age    int
	Emails []string
}

func test() {
	p := Person{
		Name:   "Alice",
		Age:    30,
		Emails: []string{"alice@example.com", "alice@gmail.com"},
	}
	fmt.Printf("Serialize result:%s \n", Serialize(p))
	var value float64
	err := Deserialize(Serialize(1.3655), &value)
	fmt.Printf("Deserialize result:%f err:%s \n", value, err)
}

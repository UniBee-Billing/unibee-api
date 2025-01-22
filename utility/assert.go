package utility

import (
	"fmt"
	"strings"
)

const (
	SystemAssertPrefix = "system_assert: "
)

func Assert(check bool, message string) {
	if !check {
		panic(SystemAssertPrefix + message)
	}
}

func AssertError(err error, message string) {
	if err != nil {
		fmt.Printf("AssertError error:%s\n", err.Error())
		panic(fmt.Sprintf(SystemAssertPrefix + strings.ReplaceAll(fmt.Sprintf("%s; %s", message, err.Error()), SystemAssertPrefix, "")))
	}
}

func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

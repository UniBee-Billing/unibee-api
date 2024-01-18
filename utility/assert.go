package utility

import (
	"fmt"
)

const (
	SystemAssertPrefix = "system_assert: "
)

func Assert(check bool, message string) {
	if !check {
		panic(SystemAssertPrefix + message)
	}
}

// Try Cache 模拟，捕获内部异常，实际上框架已实现恢复机制
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

func resume() {
	defer func() {
		if exception := recover(); exception != nil {
			fmt.Printf("exception panic error:%s\n", exception)
			return
		}
	}()
}

package main

import (
	"fmt"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"go-oversea-pay/internal/cmd"
	"go-oversea-pay/internal/cmd/nacos"
	_ "go-oversea-pay/internal/logic"
	_ "go-oversea-pay/internal/packed"
	"go-oversea-pay/redismq"
	"runtime"
)

func main() {
	// https://goframe.org/pages/viewpage.action?pageId=3672072 时区处理
	//err := gtime.SetTimeZone("Asia/Shanghai")
	//if err != nil {
	//	panic(err)
	//}
	fmt.Println("Current Go version:", runtime.Version())
	nacos.Init()
	redismq.StartRedisMqConsumer()
	cmd.Main.Run(gctx.GetInitCtx())
}

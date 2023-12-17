package main

import (
	_ ""
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"go-oversea-pay/internal/cmd"
	"go-oversea-pay/internal/cmd/nacos"
	_ "go-oversea-pay/internal/logic"
	_ "go-oversea-pay/internal/packed"
	"go-oversea-pay/redismq"
)

func main() {
	nacos.Init()
	redismq.StartRedisMqConsumer()
	cmd.Main.Run(gctx.GetInitCtx())
}

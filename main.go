package main

import _ "go-oversea-pay/time"

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"go-oversea-pay/internal/cmd"
	_system_config "go-oversea-pay/internal/cmd/config"
	_ "go-oversea-pay/internal/consumer"
	_ "go-oversea-pay/internal/logic"
	"go-oversea-pay/redismq"
)

func main() {
	_system_config.Init()
	redismq.StartRedisMqConsumer()
	cmd.Main.Run(gctx.GetInitCtx())
}

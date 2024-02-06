package main

import _ "unibee-api/time"

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"unibee-api/internal/cmd"
	_system_config "unibee-api/internal/cmd/config"
	_ "unibee-api/internal/consumer"
	_ "unibee-api/internal/logic"
	"unibee-api/redismq"
)

func main() {
	_system_config.Init()
	redismq.StartRedisMqConsumer()
	cmd.Main.Run(gctx.GetInitCtx())
}

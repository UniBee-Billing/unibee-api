package main

import (
	"unibee/internal/cmd/config"
	_ "unibee/time"
)

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/os/gctx"
	"unibee/internal/cmd"
	_ "unibee/internal/consumer"
	_ "unibee/internal/driver/pgsql"
	_ "unibee/internal/logic"
)

func main() {
	config.Init()
	cmd.Main.Run(gctx.GetInitCtx())
}

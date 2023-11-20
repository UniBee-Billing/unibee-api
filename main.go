package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"go-oversea-pay/internal/cmd/nacos"

	_ "go-oversea-pay/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"go-oversea-pay/internal/cmd"
)

func main() {
	nacos.Init()
	cmd.Main.Run(gctx.GetInitCtx())
}

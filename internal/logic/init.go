package logic

import (
	"context"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"time"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/email"
	"unibee/internal/logic/merchant"
)

func printStandaloneInitPanic(ctx context.Context, err error) {
	if err != nil {
		g.Log().Errorf(ctx, "StandaloneInit panic error:%s", err.Error())
	} else {
		g.Log().Errorf(ctx, "StandaloneInit panic error:%s", err)
	}
}

func StandaloneInit(ctx context.Context) {
	if config.GetConfigInstance().Mode != "cloud" {
		go func() {
			backgroundCtx := context.Background()
			var err error
			defer func() {
				if exception := recover(); exception != nil {
					if v, ok := exception.(error); ok && gerror.HasStack(v) {
						err = v
					} else {
						err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
					}
					printStandaloneInitPanic(backgroundCtx, err)
					return
				}
			}()
			time.Sleep(10 * time.Second)
			merchant.StandAloneInit(backgroundCtx)
			email.StandAloneInit(backgroundCtx)
		}()
	}
}

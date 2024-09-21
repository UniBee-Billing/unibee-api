package log

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
)

func PrintPanic(ctx context.Context, err error) {
	if err != nil {
		g.Log().Errorf(ctx, "WebhookSend panic error:%s", err.Error())
	} else {
		g.Log().Errorf(ctx, "WebhookSend panic error:%s", err)
	}
}

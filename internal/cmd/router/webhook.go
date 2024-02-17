package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	_webhook "unibee-api/internal/controller/webhook"
)

func Webhook(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			_webhook.NewSetup(),
		)
	})
}

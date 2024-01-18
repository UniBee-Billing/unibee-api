package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/system"
)

func SystemSubscription(ctx context.Context, group *ghttp.RouterGroup) {
	// profile 库相关接口
	group.Group("/subscription", func(group *ghttp.RouterGroup) {
		group.Bind(
			system.NewSubscription(),
		)
	})
}

func SystemInvoice(ctx context.Context, group *ghttp.RouterGroup) {
	// profile 库相关接口
	group.Group("/invoice", func(group *ghttp.RouterGroup) {
		group.Bind(
			system.NewInvoice(),
		)
	})
}

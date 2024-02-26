package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/controller/system"
)

func SystemSubscription(ctx context.Context, group *ghttp.RouterGroup) {
	// profile
	group.Group("/subscription", func(group *ghttp.RouterGroup) {
		group.Bind(
			system.NewSubscription(),
		)
	})
}

func SystemInvoice(ctx context.Context, group *ghttp.RouterGroup) {
	// profile
	group.Group("/invoice", func(group *ghttp.RouterGroup) {
		group.Bind(
			system.NewInvoice(),
		)
	})
}

func SystemPayment(ctx context.Context, group *ghttp.RouterGroup) {
	// profile
	group.Group("/payment", func(group *ghttp.RouterGroup) {
		group.Bind(
			system.NewPayment(),
		)
	})
}

func SystemRefund(ctx context.Context, group *ghttp.RouterGroup) {
	// refund
	group.Group("/refund", func(group *ghttp.RouterGroup) {
		group.Bind(
			system.NewRefund(),
		)
	})
}

func SystemMerchantInformation(ctx context.Context, group *ghttp.RouterGroup) {
	// Information
	group.Group("/merchant", func(group *ghttp.RouterGroup) {
		group.Bind(
			system.NewInformation(),
		)
	})
}

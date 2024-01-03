package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/subscription"
	"go-oversea-pay/internal/controller/subscription_plan_merchant"
	"go-oversea-pay/internal/controller/webhook"
)

// SubscriptionPlanAdmin 订阅类
func SubscriptionPlanAdmin(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service 库相关接口
	group.Group("/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			subscription_plan_merchant.NewV1(), //Plan Admin接口-Merchant Portal使用
		)
	})
}

func Subscription(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service 库相关接口
	group.Group("/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			subscription.NewV1(), //Subscription接口-User Portal使用
		)
	})
}

func Webhook(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service 库相关接口
	group.Group("/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			webhook.NewV1(), //Webhook接口-Go Server使用
		)
	})
}

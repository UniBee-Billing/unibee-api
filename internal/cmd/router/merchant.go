package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/merchant"
)

// MerchantPlan 订阅类
func MerchantPlan(ctx context.Context, group *ghttp.RouterGroup) {
	// plan 库相关接口
	group.Group("/plan", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewPlan(), //MerchantPlan Admin接口-Merchant Portal使用
		)
	})
}

func MerchantWebhook(ctx context.Context, group *ghttp.RouterGroup) {
	// auth 库相关接口
	group.Group("/auth", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewWebhook(), //Webhook接口-Go Server使用
		)
	})
}

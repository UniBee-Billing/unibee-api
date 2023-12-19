package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/subscription"
)

// Subscription 订阅类
func Subscription(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service 库相关接口
	group.Group("/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			subscription.NewV1(), //订阅平台接口
		)
	})
}

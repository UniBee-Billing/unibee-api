package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/user"
)

// UserAuth 工具类
func UserAuth(ctx context.Context, group *ghttp.RouterGroup) {
	// auth 库相关接口
	group.Group("/auth", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewAuth(),
		)
	})
}

// UserPlan 工具类
func UserPlan(ctx context.Context, group *ghttp.RouterGroup) {
	// subscription 库相关接口
	group.Group("/plan", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewPlan(),
		)
	})
}

// UserSubscription 工具类
func UserSubscription(ctx context.Context, group *ghttp.RouterGroup) {
	// subscription 库相关接口
	group.Group("/subscription", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewSubscription(),
		)
	})
}

// UserProfile 工具类
func UserProfile(ctx context.Context, group *ghttp.RouterGroup) {
	// profile 库相关接口
	group.Group("/profile", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewProfile(),
		)
	})
}

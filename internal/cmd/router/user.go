package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/user"
)

// UserAuth 工具类
func UserAuth(ctx context.Context, group *ghttp.RouterGroup) {
	// auth 库相关接口
	group.Group("/", func(group *ghttp.RouterGroup) {
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
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewProfile(),
		)
	})
}

// UserVat 工具类
func UserVat(ctx context.Context, group *ghttp.RouterGroup) {
	// vat 库相关接口
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewVat(),
		)
	})
}

// UserInvoice 工具类
func UserInvoice(ctx context.Context, group *ghttp.RouterGroup) {
	// invoice 库相关接口
	group.Group("/invoice", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewInvoice(),
		)
	})
}

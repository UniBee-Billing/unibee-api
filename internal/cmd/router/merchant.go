package router

import (
	"context"
	"unibee-api/internal/controller/merchant"

	"github.com/gogf/gf/v2/net/ghttp"
)

func MerchantUserAuth(ctx context.Context, group *ghttp.RouterGroup) {
	// auth 库相关接口
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewAuth(),
		)
	})
}

func MerchantProfile(ctx context.Context, group *ghttp.RouterGroup) {
	// profile 库相关接口
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewProfile(),
		)
	})
}

// MerchantPlan 订阅类
func MerchantPlan(ctx context.Context, group *ghttp.RouterGroup) {
	// plan 库相关接口
	group.Group("/plan", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewPlan(), //MerchantPlan Admin接口-Merchant Portal使用
		)
	})
}

// MerchantSubscrption 订阅类
func MerchantSubscrption(ctx context.Context, group *ghttp.RouterGroup) {
	// plan 库相关接口
	group.Group("/subscription", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewSubscription(), //MerchantSubscription Admin接口-Merchant Portal使用
		)
	})
}

// MerchantInvoice 发票类
func MerchantInvoice(ctx context.Context, group *ghttp.RouterGroup) {
	// plan 库相关接口
	group.Group("/invoice", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewInvoice(), //MerchantSubscription Admin接口-Merchant Portal使用
		)
	})
}

func MerchantWebhook(ctx context.Context, group *ghttp.RouterGroup) {
	// auth 库相关接口
	group.Group("/webhook", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewWebhook(), //Webhook接口-Go Server使用
		)
	})
}

func MerchantOss(ctx context.Context, group *ghttp.RouterGroup) {
	// auth 库相关接口
	group.Group("/oss", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewOss(), //Oss 文件接口-Go Server使用
		)
	})
}

func MerchantVat(ctx context.Context, group *ghttp.RouterGroup) {
	// auth 库相关接口
	group.Group("/vat", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewVat(), //Vat 文件接口-Go Server使用
		)
	})
}

func MerchantBalance(ctx context.Context, group *ghttp.RouterGroup) {
	// auth 库相关接口
	group.Group("/balance", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewBalance(), //Balance 文件接口-Go Server使用
		)
	})
}

func MerchantPayment(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/payment", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewPayment(),
		)
	})
}

func MerchantUser(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/merchant_user", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewUser(),
		)
	})
}

func MerchantInfo(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/merchant_info", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewMerchantinfo(),
		)
	})
}

func MerchantSearch(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/search", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewSearch(),
		)
	})
}

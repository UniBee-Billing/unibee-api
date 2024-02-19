package router

import (
	"context"
	"unibee-api/internal/controller/merchant"

	"github.com/gogf/gf/v2/net/ghttp"
)

func MerchantUserAuth(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewAuth(),
		)
	})
}

func MerchantProfile(ctx context.Context, group *ghttp.RouterGroup) {
	// profile
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewProfile(),
		)
	})
}

// MerchantPlan 订阅类
func MerchantPlan(ctx context.Context, group *ghttp.RouterGroup) {
	// plan
	group.Group("/plan", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewPlan(), //MerchantPlan Admin接口-Merchant Portal使用
		)
	})
}

// MerchantSubscrption 订阅类
func MerchantSubscrption(ctx context.Context, group *ghttp.RouterGroup) {
	// plan
	group.Group("/subscription", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewSubscription(), //MerchantSubscription Admin接口-Merchant Portal使用
		)
	})
}

// MerchantInvoice 发票类
func MerchantInvoice(ctx context.Context, group *ghttp.RouterGroup) {
	// plan
	group.Group("/invoice", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewInvoice(), //MerchantSubscription Admin接口-Merchant Portal使用
		)
	})
}

func MerchantGateway(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/gateway", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewGateway(), //Gateway
		)
	})
}

func MerchantOss(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/oss", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewOss(), //Oss 文件接口-Go Server使用
		)
	})
}

func MerchantVat(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/vat", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewVat(), //Vat 文件接口-Go Server使用
		)
	})
}

func MerchantBalance(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
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

func MerchantEmailTemplate(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/email_template", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewEmail(),
		)
	})
}

func MerchantWebhook(ctx context.Context, group *ghttp.RouterGroup) {
	group.Group("/merchant_webhook", func(group *ghttp.RouterGroup) {
		group.Bind(
			merchant.NewWebhook(),
		)
	})
}

package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/controller/user"
)

func UserAuth(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewAuth(),
		)
	})
}

func UserPlan(ctx context.Context, group *ghttp.RouterGroup) {
	// subscription
	group.Group("/plan", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewPlan(),
		)
	})
}

func UserSubscription(ctx context.Context, group *ghttp.RouterGroup) {
	// subscription
	group.Group("/subscription", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewSubscription(),
		)
	})
}

func UserProfile(ctx context.Context, group *ghttp.RouterGroup) {
	// profile
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewProfile(),
		)
	})
}

func UserVat(ctx context.Context, group *ghttp.RouterGroup) {
	// vat
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewVat(),
		)
	})
}

func UserInvoice(ctx context.Context, group *ghttp.RouterGroup) {
	// invoice
	group.Group("/invoice", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewInvoice(),
		)
	})
}

func UserPayment(ctx context.Context, group *ghttp.RouterGroup) {
	// invoice
	group.Group("/payment", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewPayment(),
		)
	})
}

func UserSession(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/session", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewSession(),
		)
	})
}

func UserGateway(ctx context.Context, group *ghttp.RouterGroup) {
	// auth
	group.Group("/gateway", func(group *ghttp.RouterGroup) {
		group.Bind(
			user.NewGateway(),
		)
	})
}

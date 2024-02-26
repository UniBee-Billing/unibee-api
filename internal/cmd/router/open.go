package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"unibee/internal/controller/onetime"
)

// OpenPayment 工具类的，不含业务属性的
func OpenPayment(ctx context.Context, group *ghttp.RouterGroup) {
	// payment
	group.Group("/payment", func(group *ghttp.RouterGroup) {
		group.Bind(
			onetime.NewPayment(), //开放平台接口
		)
	})
}

// OpenMocks 工具类
func OpenMocks(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service
	group.Group("/auth", func(group *ghttp.RouterGroup) {
		group.Bind(
			onetime.NewMock(), //开放平台接口
		)
	})
}

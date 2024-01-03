package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/open"
)

// OpenPayment 工具类的，不含业务属性的
func OpenPayment(ctx context.Context, group *ghttp.RouterGroup) {
	// payment 库相关接口
	group.Group("/payment", func(group *ghttp.RouterGroup) {
		group.Bind(
			open.NewPayment(), //开放平台接口
		)
	})
}

// OpenMocks 工具类
func OpenMocks(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service 库相关接口
	group.Group("/auth", func(group *ghttp.RouterGroup) {
		group.Bind(
			open.NewMock(), //开放平台接口
		)
	})
}

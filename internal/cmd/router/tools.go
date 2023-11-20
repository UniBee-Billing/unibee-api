package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/xin"
)

// Tools 工具类的，不含业务属性的
func Tools(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service 库相关接口
	group.Group("/", func(group *ghttp.RouterGroup) {
		//group.Middleware(service.Middleware().Auth)
		group.Bind(
			xin.NewV1(), //测试 xin_service 库接口
		)
	})
}

package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/out"
)

// Outs 工具类的，不含业务属性的
func Outs(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service 库相关接口
	group.Group("/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			out.NewV1(), //开放平台接口
		)
	})
}

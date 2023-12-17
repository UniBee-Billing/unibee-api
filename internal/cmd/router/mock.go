package router

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/controller/mock"
)

// Mocks 工具类
func Mocks(ctx context.Context, group *ghttp.RouterGroup) {
	// xin_service 库相关接口
	group.Group("/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			mock.NewV1(), //开放平台接口
		)
	})
}

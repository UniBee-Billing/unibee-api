package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"go-oversea-pay/internal/cmd/router"
	"go-oversea-pay/internal/cmd/swagger"
	"go-oversea-pay/internal/controller/hello"
)

var (
	Main = gcmd.Command{
		Name:  "go_oversea_pay",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.GET("/swagger", func(r *ghttp.Request) {
					r.Response.Write(swagger.SwaggerUIPageContent)
				})
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Bind(
					hello.NewV1(),
				)
			})
			s.Group("/tools", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				router.Tools(ctx, group) //工具接口
			})

			s.Run()

			return nil
		},
	}
)

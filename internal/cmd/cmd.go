package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"go-oversea-pay/internal/cmd/router"
	"go-oversea-pay/internal/cmd/swagger"
	"go-oversea-pay/internal/logic/webhooks"
	"go-oversea-pay/internal/service"
	"go-oversea-pay/utility/liberr"
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
				group.Middleware(
					service.Middleware().ResponseHandler,
					service.Middleware().PreAuth,
				)
				//group.Bind(
				//	hello.NewV1(), //测试接口
				//)
			})
			s.Group("/xin", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().ResponseHandler,
					service.Middleware().PreAuth,
				)
				router.Tools(ctx, group) //工具接口
			})

			s.Group("/out", func(group *ghttp.RouterGroup) {
				group.Middleware(
					service.Middleware().ResponseHandler,
					service.Middleware().PreAuth,
				)
				router.Outs(ctx, group) //开放平台接口
			})

			// 通道支付 Webhook 回调
			s.BindHandler("POST:/webhooks/{channelId}/notifications", webhooks.ChannelPaymentWebhookEntrance)
			// 通道支付 Redirect 回调
			s.BindHandler("GET:/redirect/{channelId}/forward", webhooks.ChannelPaymentRedirectEntrance)

			{
				_, err := g.Redis().Set(ctx, "g_check", "checked")
				liberr.ErrIsNil(ctx, err, "redis write check failure")
				value, err := g.Redis().Get(ctx, "g_check")
				liberr.ErrIsNil(ctx, err, "redis read check failure")
				g.Log().Infof(ctx, "redis check success: %s ", value.String())
			}

			s.Run()

			return nil
		},
	}
)

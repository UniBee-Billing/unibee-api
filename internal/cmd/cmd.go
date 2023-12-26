package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"go-oversea-pay/internal/cmd/router"
	"go-oversea-pay/internal/cmd/swagger"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/controller/webhooks"
	"go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/payment/outchannel"
	"go-oversea-pay/utility/liberr"
)

var (
	Main = gcmd.Command{
		Name:  "go_oversea_pay",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			s := g.Server()
			s.Group("/gooverseapay", func(group *ghttp.RouterGroup) {
				group.GET("/swagger-ui.html", func(r *ghttp.Request) {
					r.Response.Write(swagger.SwaggerUIPageContent)
				})
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreOpenApiAuth,
				)
				//group.Bind(
				//	hello.NewV1(), //测试接口
				//)
			})
			//s.Group("/gooverseapay/xin", func(group *ghttp.RouterGroup) {
			//	group.Middleware(
			//		_interface.Middleware().ResponseHandler,
			//		_interface.Middleware().PreOpenApiAuth,
			//	)
			//	router.Tools(ctx, group) //工具接口
			//})

			s.Group("/gooverseapay/out", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreOpenApiAuth,
				)
				router.Outs(ctx, group) //开放平台接口
			})

			s.Group("/gooverseapay/subscription", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreAuth,
				)
				router.Subscription(ctx, group) //订阅接口
			})

			s.Group("/gooverseapay/mock", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreAuth,
				)
				router.Mocks(ctx, group) //Out本地测试用Mock接口
			})

			s.Group("/gooverseapay/auth", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreAuth,
				)
				router.Auth(ctx, group) //Out本地测试用Mock接口
			})

			// 通道支付 Redirect 回调
			s.BindHandler("GET:/gooverseapay/payment/redirect/{channelId}/forward", webhooks.ChannelPaymentRedirectEntrance)
			// 通道支付 Webhook 回调
			s.BindHandler("POST:/gooverseapay/payment/webhooks/{channelId}/notifications", webhooks.ChannelPaymentWebhookEntrance)
			// 初始化通道 Webhook 配置
			outchannel.CheckAndSetupPayChannelWebhooks(ctx)

			{
				_, err := g.Redis().Set(ctx, "g_check", "checked")
				liberr.ErrIsNil(ctx, err, "redis write check failure")
				value, err := g.Redis().Get(ctx, "g_check")
				liberr.ErrIsNil(ctx, err, "redis read check failure")
				_, err = g.Redis().Expire(ctx, "g_check", 10)
				liberr.ErrIsNil(ctx, err, "redis write expire failure")
				g.Log().Infof(ctx, "redis check success: %s ", value.String())
				g.Log().Infof(ctx, "swagger address: http://127.0.0.1%s/gooverseapay/swagger", consts.GetConfigInstance().Server.Address)
			}

			s.Run()

			return nil
		},
	}
)

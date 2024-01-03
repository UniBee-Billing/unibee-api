package cmd

import (
	"context"
	"go-oversea-pay/internal/cmd/router"
	"go-oversea-pay/internal/cmd/swagger"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/controller"
	"go-oversea-pay/internal/controller/channel_webhook_entry"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/utility/liberr"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "go_oversea_pay",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {

			s := g.Server()
			s.Group("/"+consts.GetConfigInstance().Server.Name, func(group *ghttp.RouterGroup) {
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

			//s.Group("/"+consts.GetConfigInstance().Server.Name+"/out", func(group *ghttp.RouterGroup) {
			//	group.Middleware(
			//		_interface.Middleware().ResponseHandler,
			//		_interface.Middleware().PreOpenApiAuth,
			//	)
			//	router.Outs(ctx, group) //开放平台接口
			//})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/plan/merchant", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().TokenAuth,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreAuth,
				)
				router.SubscriptionPlanAdmin(ctx, group) //订阅Plan接口-Merchant Portal 使用
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/subscription", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
				)
				router.Subscription(ctx, group) //订阅Plan接口-Merchant Portal 使用
			})

			//s.Group("/"+consts.GetConfigInstance().Server.Name+"/mock", func(group *ghttp.RouterGroup) {
			//	group.Middleware(
			//		_interface.Middleware().ResponseHandler,
			//		_interface.Middleware().PreAuth,
			//	)
			//	router.Mocks(ctx, group) //Out本地测试用Mock接口
			//})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/auth", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreAuth,
				)
				router.Auth(ctx, group) //Out本地测试用Mock接口
			})

			s.BindHandler("GET:/health", controller.HealthCheck)

			// 通道支付 Redirect 回调
			s.BindHandler("GET:/"+consts.GetConfigInstance().Server.Name+"/payment/redirect/{channelId}/forward", channel_webhook_entry.ChannelPaymentRedirectEntrance)
			// 通道支付 Webhook 回调
			s.BindHandler("POST:/"+consts.GetConfigInstance().Server.Name+"/payment/channel_webhook_entry/{channelId}/notifications", channel_webhook_entry.ChannelPaymentWebhookEntrance)
			//// 初始化通道 Webhook 配置
			//outchannel.CheckAndSetupPayChannelWebhooks(ctx)

			{
				g.Log().Infof(ctx, "server name: %s ", consts.GetConfigInstance().Server.Name)
				g.Log().Infof(ctx, "server port: %s ", consts.GetConfigInstance().Server.Address)
				g.Log().Infof(ctx, "server domainPath: %s ", consts.GetConfigInstance().Server.DomainPath)
				_, err := g.Redis().Set(ctx, "g_check", "checked")
				liberr.ErrIsNil(ctx, err, "redis write check failure")
				value, err := g.Redis().Get(ctx, "g_check")
				liberr.ErrIsNil(ctx, err, "redis read check failure")
				_, err = g.Redis().Expire(ctx, "g_check", 10)
				liberr.ErrIsNil(ctx, err, "redis write expire failure")
				g.Log().Infof(ctx, "redis check success: %s ", value.String())
				g.Log().Infof(ctx, "swagger try address: http://127.0.0.1%s/%s/swagger-ui.html", consts.GetConfigInstance().Server.Address, consts.GetConfigInstance().Server.Name)
			}

			s.Run()

			return nil
		},
	}
)

package cmd

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"go-oversea-pay/internal/cmd/router"
	"go-oversea-pay/internal/cmd/swagger"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/controller"
	"go-oversea-pay/internal/controller/channel_webhook_entry"
	"go-oversea-pay/internal/cronjob"
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
					// _interface.Middleware().PreOpenApiAuth,
				)
				//group.Bind(
				//	hello.NewV1(), //测试接口
				//)
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/open", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreOpenApiAuth,
				)
				router.OpenPayment(ctx, group) //开放平台接口
				router.OpenMocks(ctx, group)   //Out本地测试用Mock接口
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/merchant", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().TokenMerchantAuth,
				)
				router.MerchantPlan(ctx, group)
				router.MerchantWebhook(ctx, group)
				router.MerchantProfile(ctx, group)
				router.MerchantSubscrption(ctx, group)
				router.MerchantInvoice(ctx, group)
				router.MerchantOss(ctx, group)
				router.MerchantVat(ctx, group)
				router.MerchantBalance(ctx, group)
				router.MerchantPayment(ctx, group)
				router.MerchantUser(ctx, group)
				router.MerchantSearch(ctx, group)
				router.MerchantInfo(ctx, group)
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/merchant/auth", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
				)
				router.MerchantUserAuth(ctx, group)
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/user", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().TokenUserAuth,
				)
				router.UserPlan(ctx, group)
				router.UserSubscription(ctx, group)
				router.UserInvoice(ctx, group)
				router.UserProfile(ctx, group)
				router.UserPayment(ctx, group)
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/user/auth", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
				)
				router.UserAuth(ctx, group)
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/user/vat", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
				)
				router.UserVat(ctx, group)
			})

			if !consts.GetConfigInstance().IsProd() {
				s.Group("/"+consts.GetConfigInstance().Server.Name+"/system", func(group *ghttp.RouterGroup) {
					group.Middleware(
						_interface.Middleware().CORS,
						_interface.Middleware().ResponseHandler,
					)
					router.SystemSubscription(ctx, group)
					router.SystemInvoice(ctx, group)
					router.SystemPayment(ctx, group)
					router.SystemRefund(ctx, group)
				})
			}

			s.BindHandler("GET:/health", controller.HealthCheck)

			// 通道支付 Redirect 回调
			s.BindHandler("GET:/"+consts.GetConfigInstance().Server.Name+"/payment/redirect/{channelId}/forward", channel_webhook_entry.ChannelPaymentRedirectEntrance)
			// 通道支付 MerchantWebhook 回调
			s.BindHandler("POST:/"+consts.GetConfigInstance().Server.Name+"/payment/channel_webhook_entry/{channelId}/notifications", channel_webhook_entry.ChannelPaymentWebhookEntrance)
			//// 初始化通道 MerchantWebhook 配置
			//channel.CheckAndSetupPayChannelWebhooks(ctx)

			{
				g.Log().Infof(ctx, "Server name: %s ", consts.GetConfigInstance().Server.Name)
				g.Log().Infof(ctx, "Server port: %s ", consts.GetConfigInstance().Server.Address)
				//g.Log().Infof(ctx, "Server TimeZone: %d ", time.z)
				g.Log().Infof(ctx, "Server TimeStamp: %d ", gtime.Now().Timestamp())
				g.Log().Infof(ctx, "Server Time: %s ", gtime.Now().Layout("2006-01-02 15:04:05"))
				g.Log().Infof(ctx, "Server domainPath: %s ", consts.GetConfigInstance().Server.DomainPath)
				_, err := g.Redis().Set(ctx, "g_check", "checked")
				liberr.ErrIsNil(ctx, err, "Redis write check failure")
				value, err := g.Redis().Get(ctx, "g_check")
				liberr.ErrIsNil(ctx, err, "Redis read check failure")
				_, err = g.Redis().Expire(ctx, "g_check", 10)
				liberr.ErrIsNil(ctx, err, "Redis write expire failure")
				g.Log().Infof(ctx, "Redis check success: %s ", value.String())
				g.Log().Infof(ctx, "Swagger try address: http://127.0.0.1%s/%s/swagger-ui.html", consts.GetConfigInstance().Server.Address, consts.GetConfigInstance().Server.Name)
				if !consts.GetConfigInstance().IsServerDev() && !consts.GetConfigInstance().IsLocal() {
					//g.Log().SetLevel(glog.LEVEL_ALL ^ glog.LEVEL_DEBU) // remote debug log
					_ = g.Log().SetLevelStr("info") // remote debug log, DEBU < INFO < NOTI < WARN < ERRO < CRIT
				}
				cronjob.StartCronJobs()
			}

			s.Run()

			return nil
		},
	}
)

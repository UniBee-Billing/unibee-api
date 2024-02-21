package cmd

import (
	"context"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
	"unibee-api/internal/cmd/router"
	"unibee-api/internal/cmd/swagger"
	"unibee-api/internal/consts"
	"unibee-api/internal/controller"
	"unibee-api/internal/controller/gateway_webhook_entry"
	"unibee-api/internal/cronjob"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/session"
	"unibee-api/utility"
	"unibee-api/utility/liberr"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "UniBee Api",
		Usage: "main",
		Brief: "start server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.Group("/"+consts.GetConfigInstance().Server.Name, func(group *ghttp.RouterGroup) {
				group.GET("/swagger-ui.html", func(r *ghttp.Request) {
					r.Response.Write(swagger.LatestSwaggerUIPageContent)
				})
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					// _interface.Middleware().PreOpenApiAuth,
				)
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/session", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreOpenApiAuth,
				)
				router.UserSession(ctx, group)
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/open", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().PreOpenApiAuth,
				)
				router.OpenPayment(ctx, group)
				router.OpenMocks(ctx, group)
			})

			s.Group("/"+consts.GetConfigInstance().Server.Name+"/merchant", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().TokenMerchantAuth,
				)
				router.MerchantPlan(ctx, group)
				router.MerchantGateway(ctx, group)
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
				router.MerchantEmailTemplate(ctx, group)
				router.MerchantWebhook(ctx, group)
				router.MerchantMetric(ctx, group)
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
					_interface.Middleware().TokenUserAuth,
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
					router.SystemMerchantInformation(ctx, group)
				})
			}

			s.BindHandler("GET:/health", controller.HealthCheck)

			// Session Redirect
			s.BindHandler("GET:/"+consts.GetConfigInstance().Server.Name+"/session/redirect/{session}/forward", session.UserSessionRedirectEntrance)

			// Gateway Redirect
			s.BindHandler("GET:/"+consts.GetConfigInstance().Server.Name+"/payment/redirect/{gatewayId}/forward", gateway_webhook_entry.GatewayRedirectEntrance)
			// Gateway Webhook
			s.BindHandler("POST:/"+consts.GetConfigInstance().Server.Name+"/payment/gateway_webhook_entry/{gatewayId}/notifications", gateway_webhook_entry.GatewayWebhookEntrance)

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
					//g.Log().SetLevel(glog.LEVEL_ALL ^ glog.LEVEL_DEBU) // use bit migrate remove debug log
					_ = g.Log().SetLevelStr("info") // remove debug log, DEBU < INFO < NOTI < WARN < ERRO < CRIT
				}
				cronjob.StartCronJobs()
			}
			{
				g.Log().Infof(ctx, "TimeZone:%s", utility.MarshalToJsonString(time.Local))
			}

			s.Run()

			return nil
		},
	}
)

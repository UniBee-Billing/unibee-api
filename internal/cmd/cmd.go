package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gtime"
	"time"
	"unibee/internal/cmd/swagger"
	"unibee/internal/consts"
	"unibee/internal/consumer/websocket"
	"unibee/internal/controller"
	"unibee/internal/controller/gateway_webhook_entry"
	"unibee/internal/controller/link/invoice"
	"unibee/internal/controller/link/payment"
	"unibee/internal/controller/merchant"
	"unibee/internal/controller/system"
	"unibee/internal/controller/user"
	"unibee/internal/cronjob"
	_interface "unibee/internal/interface"
	"unibee/utility"
	"unibee/utility/liberr"

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

			_, err = g.Cfg().Get(ctx, "database")
			if err != nil {
				return err
			}

			openapi := s.GetOpenApi()
			openapi.Info.Description = "This is UniBee api server, For this sample, you can use the api key `EUXAgwv3Vcr1PFWt2SgBumMHXn3ImBqM` to test the authorization filters"
			openapi.Info.Title = "OpenAPI UniBee"
			openapi.Security = &goai.SecurityRequirements{map[string][]string{}}
			openapi.Info.License = &goai.License{
				Name: "Apache-2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			}
			openapi.Config.CommonResponse = _interface.JsonRes{}
			openapi.Config.CommonResponseDataField = `Data`
			//if consts.GetConfigInstance().IsLocal() || consts.GetConfigInstance().IsServerDev() {
			//	openapi.Servers = &goai.Servers{goai.Server{
			//		URL:         "http://127.0.0.1" + consts.GetConfigInstance().Server.Address,
			//		Description: consts.GetConfigInstance().Env,
			//	}, goai.Server{
			//		URL:         "http://api.unibee.top",
			//		Description: "stage",
			//	}}
			//} else {
			//	openapi.Servers = &goai.Servers{goai.Server{
			//		URL:         consts.GetConfigInstance().Server.DomainPath,
			//		Description: consts.GetConfigInstance().Env,
			//	}}
			//}

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.GET("/swagger-ui.html", func(r *ghttp.Request) {
					r.Response.Write(swagger.V3SwaggerUIPageContent)
				})
				group.GET("/api.sdk.generator.json", swagger.MerchantPortalAndSDKGeneratorSpecJson)
				group.GET("/api.user.portal.generator.json", swagger.UserPortalGeneratorSpecJson)
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
				)
			})

			s.Group("/merchant", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().TokenAuth,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewProfile(),
					)
				})
				group.Group("/plan", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewPlan(),
					)
				})
				group.Group("/gateway", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewGateway(),
					)
				})
				group.Group("/member", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewMember(),
					)
				})
				group.Group("/subscription", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewSubscription(),
					)
				})
				group.Group("/invoice", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewInvoice(),
					)
				})
				group.Group("/oss", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewOss(),
					)
				})
				group.Group("/vat", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewVat(),
					)
				})
				group.Group("/balance", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewBalance(),
					)
				})
				group.Group("/payment", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewPayment(),
					)
				})
				group.Group("/user", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewUser(),
					)
				})
				group.Group("/search", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewSearch(),
					)
				})
				group.Group("/email", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewEmail(),
					)
				})
				group.Group("/webhook", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewWebhook(),
					)
				})
				group.Group("/metric", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewMetric(),
					)
				})
				group.Group("/session", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewSession(),
					)
				})
			})

			s.Group("/user", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().TokenAuth,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewProfile(),
					)
				})
				group.Group("/plan", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewPlan(),
					)
				})
				group.Group("/subscription", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewSubscription(),
					)
				})
				group.Group("/invoice", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewInvoice(),
					)
				})
				group.Group("/payment", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewPayment(),
					)
				})
				group.Group("/gateway", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewGateway(),
					)
				})
				group.Group("/merchant", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewMerchantinfo(),
					)
				})
				group.Group("/vat", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewVat(),
					)
				})
			})

			s.Group("/merchant/auth", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewAuth(),
					)
				})
			})

			s.Group("/user/auth", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().UserPortalPreAuth,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewAuth(),
					)
				})
			})

			s.Group("/system", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().OpenApiDetach,
				)
				group.Group("/information", func(group *ghttp.RouterGroup) {
					group.Bind(
						system.NewInformation(),
					)
				})
				if !consts.GetConfigInstance().IsProd() {
					group.Group("/subscription", func(group *ghttp.RouterGroup) {
						group.Bind(
							system.NewSubscription(),
						)
					})
					group.Group("/invoice", func(group *ghttp.RouterGroup) {
						group.Bind(
							system.NewInvoice(),
						)
					})
					group.Group("/payment", func(group *ghttp.RouterGroup) {
						group.Bind(
							system.NewPayment(),
						)
					})
					group.Group("/refund", func(group *ghttp.RouterGroup) {
						group.Bind(
							system.NewRefund(),
						)
					})
				}
			})

			s.BindHandler("GET:/health", controller.HealthCheck)

			// Invoice Link
			s.BindHandler("GET:/in/{invoiceId}", invoice.InvoiceLinkEntry)
			s.BindHandler("GET:/pay/{paymentId}", payment.PaymentLinkEntry)
			// Gateway Redirect
			s.BindHandler("GET:/payment/redirect/{gatewayId}/forward", gateway_webhook_entry.GatewayRedirectEntrance)
			// Gateway Webhook
			s.BindHandler("POST:/payment/gateway_webhook_entry/{gatewayId}/notifications", gateway_webhook_entry.GatewayWebhookEntrance)
			// Merchant Websocket
			s.BindHandler("/merchant_ws/{merchantApiKey}", websocket.MerchantWebSocketMessageEntry)

			{
				g.Log().Infof(ctx, "Server port: %s ", consts.GetConfigInstance().Server.Address)
				g.Log().Infof(ctx, "Server domainPath: %s ", consts.GetConfigInstance().Server.DomainPath)
				g.Log().Infof(ctx, "Server TimeStamp: %d ", gtime.Now().Timestamp())
				g.Log().Infof(ctx, "Server Time: %s ", gtime.Now().Layout("2006-01-02 15:04:05"))
				_, err := g.Redis().Set(ctx, "g_check", "checked")
				liberr.ErrIsNil(ctx, err, "Redis write check failure")
				value, err := g.Redis().Get(ctx, "g_check")
				liberr.ErrIsNil(ctx, err, "Redis read check failure")
				_, err = g.Redis().Expire(ctx, "g_check", 10)
				liberr.ErrIsNil(ctx, err, "Redis write expire failure")
				g.Log().Infof(ctx, "Redis check success: %s ", value.String())
				g.Log().Infof(ctx, "Swagger try address: http://127.0.0.1%s/swagger-ui.html", consts.GetConfigInstance().Server.Address)
				if !consts.GetConfigInstance().IsServerDev() && !consts.GetConfigInstance().IsLocal() {
					_ = g.Log().SetLevelStr("info") // remove debug log, DEBU < INFO < NOTI < WARN < ERRO < CRIT
				}
				cronjob.StartCronJobs()
			}
			{
				g.Log().Infof(ctx, "TimeZone:%s", utility.MarshalToJsonString(time.Local))
			}
			{
				fmt.Println(utility.MarshalToJsonString(s.GetOpenApi()))
			}

			s.Run()

			return nil
		},
	}
)

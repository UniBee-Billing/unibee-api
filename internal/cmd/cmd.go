package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/goai"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"time"
	"unibee/internal/cmd/config"
	"unibee/internal/cmd/swagger"
	"unibee/internal/consumer/websocket"
	"unibee/internal/controller"
	"unibee/internal/controller/gateway_webhook_entry"
	"unibee/internal/controller/link/_import"
	"unibee/internal/controller/link/export"
	"unibee/internal/controller/link/invoice"
	"unibee/internal/controller/link/oss"
	"unibee/internal/controller/link/payment"
	"unibee/internal/controller/merchant"
	"unibee/internal/controller/system"
	"unibee/internal/controller/user"
	"unibee/internal/cronjob"
	_interface "unibee/internal/interface"
	"unibee/internal/logic"
	"unibee/internal/logic/dbupgrade"
	"unibee/internal/logic/email/gateway"
	"unibee/internal/logic/gateway/webhook"
	"unibee/internal/logic/member"
	merchant2 "unibee/internal/logic/merchant"
	"unibee/internal/query"
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
			openapi := s.GetOpenApi()
			openapi.Info.Description = "UniBee Api Server"
			openapi.Info.Title = "OpenAPI UniBee"
			openapi.Info.License = &goai.License{
				Name: "Apache-2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			}
			openapi.Info.Version = utility.ReadBuildVersionInfo(ctx)
			openapi.Config.CommonResponse = _interface.JsonRes{}
			openapi.Config.CommonResponseDataField = `Data`
			//https://github.com/gogf/gf/issues/3747
			openapi.Security = &goai.SecurityRequirements{{"Authorization": []string{}}}
			openapi.Components.SecuritySchemes = goai.SecuritySchemes{
				"Authorization": goai.SecuritySchemeRef{
					Value: &goai.SecurityScheme{
						Type:         "http",
						Scheme:       "bearer",
						BearerFormat: "JWT",
					},
				},
			}
			openapi.Servers = &goai.Servers{
				{URL: config.GetConfigInstance().Server.DomainPath},
				{URL: fmt.Sprintf("http://127.0.0.1%s", config.GetConfigInstance().Server.Address)},
			}
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.GET("/swagger-ui.html", func(r *ghttp.Request) {
					r.Response.Write(swagger.LatestSwaggerUIPageContent)
				})
				group.GET("/swaggerv3-ui.html", func(r *ghttp.Request) {
					r.Response.Write(swagger.LatestSwaggerUIPageContent)
				})
				group.GET("/swagger", func(r *ghttp.Request) {
					r.Response.Write(swagger.RedoclyContent)
				})
				group.GET("/api.sdk.generator.json", swagger.MerchantPortalAndSDKGeneratorSpecJson)
				group.GET("/api.sdk.generator.yaml", swagger.MerchantPortalAndSDKGeneratorSpecYaml)
				group.GET("/api.user.portal.generator.json", swagger.UserPortalGeneratorSpecJson)
				group.GET("/api.system.generator.json", swagger.SystemGeneratorSpecJson)
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
				)
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

			s.Group("/merchant", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().MerchantHandler,
				)
				group.Group("/product", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewProduct(),
					)
				})
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
				group.Group("/role", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewRole(),
					)
				})
				group.Group("/discount", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewDiscount(),
					)
				})
				group.Group("/credit", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewCredit(),
					)
				})
				group.Group("/task", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewTask(),
					)
				})
				group.Group("/track", func(group *ghttp.RouterGroup) {
					group.Bind(
						merchant.NewTrack(),
					)
				})
			})

			s.Group("/user", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().UserPortalApiHandler,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewProfile(),
					)
				})
				group.Group("/product", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewProduct(),
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
				group.Group("/metric", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewMetric(),
					)
				})
			})

			s.Group("/user/vat", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().UserPortalMerchantRouterHandler,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Bind(
						user.NewVat(),
					)
				})
			})

			s.Group("/user/auth", func(group *ghttp.RouterGroup) {
				group.Middleware(
					_interface.Middleware().CORS,
					_interface.Middleware().ResponseHandler,
					_interface.Middleware().UserPortalMerchantRouterHandler,
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
				)
				group.Group("/plan", func(group *ghttp.RouterGroup) {
					group.Bind(
						system.NewPlan(),
					)
				})
				group.Group("/information", func(group *ghttp.RouterGroup) {
					group.Bind(
						system.NewInformation(),
					)
				})
				group.Group("/subscription", func(group *ghttp.RouterGroup) {
					group.Bind(
						system.NewSubscription(),
					)
				})
				group.Group("/payment", func(group *ghttp.RouterGroup) {
					group.Bind(
						system.NewPayment(),
					)
				})
				group.Group("/invoice", func(group *ghttp.RouterGroup) {
					group.Bind(
						system.NewInvoice(),
					)
				})
				if !config.GetConfigInstance().IsProd() {
					group.Group("/refund", func(group *ghttp.RouterGroup) {
						group.Bind(
							system.NewRefund(),
						)
					})
					group.Group("/auth", func(group *ghttp.RouterGroup) {
						group.Bind(
							system.NewAuth(),
						)
					})
				}
			})

			s.BindHandler("/health", controller.HealthCheck)
			s.BindHandler("/version", controller.Version)

			// Invoice Link
			s.BindHandler("GET:/in/{invoiceId}", invoice.LinkEntry)
			s.BindHandler("GET:/in/pdf/{invoiceId}", invoice.LinkPdfEntry)
			s.BindHandler("GET:/oss/file/{filename}", oss.FileEntry)
			s.BindHandler("GET:/export/{taskId}", export.LinkExportEntry)
			s.BindHandler("GET:/import/template/{task}", _import.LinkImportTemplateEntry)
			s.BindHandler("GET:/pay/{paymentId}", payment.LinkEntry)
			// Gateway Payment Redirect
			s.BindHandler("GET:/payment/redirect/{gatewayId}/forward", gateway_webhook_entry.GatewayRedirectEntrance)
			// Gateway Payment Method Redirect
			s.BindHandler("GET:/payment/method/redirect/{gatewayId}/forward", gateway_webhook_entry.GatewayPaymentMethodRedirectEntrance)
			// Gateway Webhook
			//s.BindHandler("POST:/payment/gateway_webhook_entry/{gatewayId}/notifications", gateway_webhook_entry.GatewayWebhookEntrance)
			s.BindHandler("/payment/gateway_webhook_entry/{gatewayId}/notifications", gateway_webhook_entry.GatewayWebhookEntrance)
			// Merchant Websocket
			s.BindHandler("/merchant_ws/{merchantApiKey}", websocket.MerchantWebSocketMessageEntry)

			{
				//db check
				dbupgrade.StandAloneInit(ctx)
				_, err = query.GetMerchantList(ctx)
				liberr.ErrIsNil(ctx, err, "DB Not Ready")
				g.Log().Infof(ctx, "TimeZone:%s", utility.MarshalToJsonString(time.Local))
				g.Log().Infof(ctx, "Server build version: %s ", g.Server().GetOpenApi().Info.Version)
				g.Log().Infof(ctx, "Server port: %s ", config.GetConfigInstance().Server.Address)
				g.Log().Infof(ctx, "Server domainPath: %s ", config.GetConfigInstance().Server.DomainPath)
				g.Log().Infof(ctx, "Server TimeStamp: %d ", gtime.Now().Timestamp())
				g.Log().Infof(ctx, "Server Time: %s ", gtime.Now().Layout("2006-01-02 15:04:05"))
				_, err = g.Redis().Set(ctx, "g_check", "checked")
				liberr.ErrIsNil(ctx, err, "Redis write check failure")
				value, err := g.Redis().Get(ctx, "g_check")
				liberr.ErrIsNil(ctx, err, "Redis read check failure")
				_, err = g.Redis().Expire(ctx, "g_check", 10)
				liberr.ErrIsNil(ctx, err, "Redis write expire failure")
				g.Log().Infof(ctx, "Redis check success: %s ", value.String())
				g.Log().Infof(ctx, "Public IP: %s ", utility.GetPublicIP())
				g.Log().Infof(ctx, "Public IP: %s ", utility.GetPublicIP())
				g.Log().Infof(ctx, "Redocly address: http://127.0.0.1%s/swagger", config.GetConfigInstance().Server.Address)
				g.Log().Infof(ctx, "SwaggerV3 address: http://127.0.0.1%s/swagger-ui.html", config.GetConfigInstance().Server.Address)
				if !config.GetConfigInstance().IsServerDev() && !config.GetConfigInstance().IsLocal() {
					_ = g.Log().SetLevelStr("info") // remove debug log, DEBU < INFO < NOTI < WARN < ERRO < CRIT
				}
			}
			{
				cronjob.StartCronJobs()
			}
			{
				fmt.Println(utility.MarshalToJsonString(s.GetOpenApi()))
			}

			{
				//logic init
				logic.StandaloneInit(ctx)
				//SetupAllWebhooks
				webhook.SetupAllWebhooksBackground()
			}

			{
				redismq.RegisterRedisMqConfig(&redismq.RedisMqConfig{
					Addr:     config.GetConfigInstance().RedisConfig.Default.Address,
					Password: config.GetConfigInstance().RedisConfig.Default.Pass,
					Database: config.GetConfigInstance().RedisConfig.Default.DB,
					Group:    "GID_UniBee_Recurring",
				})
				redismq.StartRedisMqConsumer()
			}
			{
				merchant2.ReloadAllMerchantsCacheForSDKAuthBackground()
				member.ReloadAllMembersCacheForSDKAuthBackground()
			}

			{
				_, err := gateway.ReadEmailHtmlTemplate()
				if err != nil {
					g.Log().Errorf(ctx, "ReadEmailHtmlTemplate error:%s\n", err.Error())
				} else {
					g.Log().Infof(ctx, "ReadEmailHtmlTemplate success\n")
				}
			}

			s.Run()

			return nil
		},
	}
)

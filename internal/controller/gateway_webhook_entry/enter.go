package gateway_webhook_entry

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/url"
	"strconv"
	"strings"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/gateway/util"
	"unibee/internal/logic/gateway/webhook"
	session2 "unibee/internal/logic/session"
	"unibee/internal/query"
	"unibee/utility"
)

func GatewayWebhookEntrance(r *ghttp.Request) {
	gatewayId := r.Get("gatewayId").String()
	gatewayIdInt, err := strconv.Atoi(gatewayId)
	if err != nil {
		g.Log().Errorf(r.Context(), "GatewayWebhookEntrance panic url: %s gatewayId: %s err:%s", r.GetUrl(), gatewayId, err)
		return
	}
	gateway := util.GetGatewayById(r.Context(), uint64(gatewayIdInt))
	webhook.GetGatewayWebhookServiceProvider(r.Context(), uint64(gatewayIdInt)).GatewayWebhook(r, gateway)
}

func GatewayRedirectEntrance(r *ghttp.Request) {
	gatewayId := r.Get("gatewayId").String()

	gatewayIdInt, err := strconv.Atoi(gatewayId)
	if err != nil {
		g.Log().Errorf(r.Context(), "GatewayRedirectEntrance panic url:%s gatewayId: %s err:%s", r.GetUrl(), gatewayId, err)
		return
	}
	gateway := util.GetGatewayById(r.Context(), uint64(gatewayIdInt))
	redirect, err := webhook.GetGatewayWebhookServiceProvider(r.Context(), uint64(gatewayIdInt)).GatewayRedirect(r, gateway)
	if err != nil {
		r.Response.Writeln(fmt.Sprintf("%v", err))
		return
	}
	var targetUrl = ""
	if len(redirect.ReturnUrl) > 0 {
		if !redirect.Success {
			targetUrl = fmt.Sprintf("%s", redirect.ReturnUrl)
		} else if !strings.Contains(redirect.ReturnUrl, "?") {
			targetUrl = fmt.Sprintf("%s?%s", redirect.ReturnUrl, redirect.QueryPath)
		} else {
			targetUrl = fmt.Sprintf("%s&%s", redirect.ReturnUrl, redirect.QueryPath)
		}
	} else {
		merchant := query.GetMerchantById(r.Context(), gateway.MerchantId)
		if merchant != nil && len(merchant.Host) > 0 {
			if strings.HasPrefix(merchant.Host, "http") {
				targetUrl = merchant.Host
			} else {
				targetUrl = fmt.Sprintf("http://%s", merchant.Host)
			}
		} else {
			//r.Response.Writeln(utility.FormatToJsonString(redirect))
		}
	}
	if redirect.Payment != nil && len(config.GetConfigInstance().Server.HostedPagePath) > 0 && r.Get("success").Bool() {
		cancelUrl := util.GetPaymentRedirectUrl(r.Context(), redirect.Payment, "false")
		_, userSession, err := session2.NewUserSession(r.Context(), redirect.Payment.MerchantId, redirect.Payment.UserId, targetUrl, cancelUrl)
		if err == nil && len(userSession) > 0 {
			targetUrl = fmt.Sprintf("%s/payment_checker?merchantId=%d&paymentId=%s&session=%s&env=%s", config.GetConfigInstance().Server.HostedPagePath, redirect.Payment.MerchantId, redirect.Payment.PaymentId, userSession, config.GetConfigInstance().Env)
		}
	}
	if len(targetUrl) > 0 {
		r.Response.RedirectTo(targetUrl)
	} else {
		r.Response.Writeln(utility.FormatToJsonString(redirect))
	}
}

func GatewayPaymentMethodRedirectEntrance(r *ghttp.Request) {
	gatewayId := r.Get("gatewayId").String()
	gatewayIdInt, err := strconv.Atoi(gatewayId)
	if err != nil {
		g.Log().Errorf(r.Context(), "GatewayRedirectEntrance panic url:%s gatewayId: %s err:%s", r.GetUrl(), gatewayId, err)
		return
	}
	gateway := util.GetGatewayById(r.Context(), uint64(gatewayIdInt))
	utility.Assert(gateway != nil, "gateway invalid")
	err = webhook.GetGatewayWebhookServiceProvider(r.Context(), uint64(gatewayIdInt)).GatewayNewPaymentMethodRedirect(r, gateway)
	utility.AssertError(err, "System Error")
	redirectUrl, _ := url.QueryUnescape(r.Get("redirectUrl").String())
	success := r.Get("success").Bool()
	if len(redirectUrl) > 0 {
		if !strings.Contains(redirectUrl, "?") {
			r.Response.RedirectTo(fmt.Sprintf("%s?success=%v", redirectUrl, success))
		} else {
			r.Response.RedirectTo(fmt.Sprintf("%s&success=%v", redirectUrl, success))
		}
	} else {
		r.Response.Writeln(success)
	}
}

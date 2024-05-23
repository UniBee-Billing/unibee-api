package gateway_webhook_entry

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/url"
	"strconv"
	"strings"
	"unibee/internal/logic/gateway/util"
	"unibee/internal/logic/gateway/webhook"
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
	if len(redirect.ReturnUrl) > 0 {
		if !strings.Contains(redirect.ReturnUrl, "?") {
			r.Response.RedirectTo(fmt.Sprintf("%s?%s", redirect.ReturnUrl, redirect.QueryPath))
		} else {
			r.Response.RedirectTo(fmt.Sprintf("%s&%s", redirect.ReturnUrl, redirect.QueryPath))
		}
	} else {
		merchant := query.GetMerchantById(r.Context(), gateway.MerchantId)
		if merchant != nil && len(merchant.Host) > 0 {
			if strings.HasPrefix(merchant.Host, "http") {
				r.Response.RedirectTo(merchant.Host)
			} else {
				r.Response.RedirectTo(fmt.Sprintf("http://%s", merchant.Host))
			}
		} else {
			r.Response.Writeln(utility.FormatToJsonString(redirect))
		}
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

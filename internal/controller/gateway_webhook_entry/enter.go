package gateway_webhook_entry

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/gateway/util"
	"go-oversea-pay/internal/logic/gateway/webhook"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
)

func GatewayWebhookEntrance(r *ghttp.Request) {
	gatewayId := r.Get("gatewayId").String()
	gatewayIdInt, err := strconv.Atoi(gatewayId)
	if err != nil {
		g.Log().Errorf(r.Context(), "GatewayWebhookEntrance panic gatewayId: %s err:%s", r.GetUrl(), gatewayId, err)
		return
	}
	gateway := util.GetGatewayById(r.Context(), int64(gatewayIdInt))
	webhook.GetGatewayWebhookServiceProvider(r.Context(), int64(gatewayIdInt)).GatewayWebhook(r, gateway)
}

func GatewayRedirectEntrance(r *ghttp.Request) {
	gatewayId := r.Get("gatewayId").String()
	gatewayIdInt, err := strconv.Atoi(gatewayId)
	if err != nil {
		g.Log().Errorf(r.Context(), "GatewayRedirectEntrance panic gatewayId: %s err:%s", r.GetUrl(), gatewayId, err)
		return
	}
	gateway := util.GetGatewayById(r.Context(), int64(gatewayIdInt))
	redirect, err := webhook.GetGatewayWebhookServiceProvider(r.Context(), int64(gatewayIdInt)).GatewayRedirect(r, gateway)
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
		r.Response.Writeln(utility.FormatToJsonString(redirect))
	}
}
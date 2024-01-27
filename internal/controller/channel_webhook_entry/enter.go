package channel_webhook_entry

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/channel/out"
	"go-oversea-pay/internal/logic/channel/util"
	"go-oversea-pay/utility"
	"strconv"
	"strings"
)

func ChannelPaymentWebhookEntrance(r *ghttp.Request) {
	channelId := r.Get("channelId").String()
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		g.Log().Errorf(r.Context(), "ChannelPaymentWebhookEntrance panic channelId: %s err:%s", r.GetUrl(), channelId, err)
		return
	}
	payChannel := util.GetOverseaPayChannel(r.Context(), int64(channelIdInt))
	out.GetPayChannelServiceProvider(r.Context(), int64(channelIdInt)).DoRemoteChannelWebhook(r, payChannel)
}

func ChannelPaymentRedirectEntrance(r *ghttp.Request) {
	channelId := r.Get("channelId").String()
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		g.Log().Errorf(r.Context(), "ChannelPaymentRedirectEntrance panic channelId: %s err:%s", r.GetUrl(), channelId, err)
		return
	}
	payChannel := util.GetOverseaPayChannel(r.Context(), int64(channelIdInt))
	redirect, err := out.GetPayChannelServiceProvider(r.Context(), int64(channelIdInt)).DoRemoteChannelRedirect(r, payChannel)
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

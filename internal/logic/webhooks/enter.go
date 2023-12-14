package webhooks

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/payment/outchannel"
	"strconv"
)

func ChannelPaymentWebhookEntrance(r *ghttp.Request) {
	channelId := r.Get("channelId").String()
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		g.Log().Errorf(r.Context(), "ChannelPaymentWebhookEntrance panic channelId: %s err:%s", r.GetUrl(), channelId, err)
		return
	}
	//outchannel := util.GetOverseaPayChannel(r.Context(), uint64(channelIdInt))
	outchannel.GetPayChannelServiceProvider(channelIdInt).DoRemoteChannelWebhook(r)
}

func ChannelPaymentRedirectEntrance(r *ghttp.Request) {
	channelId := r.Get("channelId").String()
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		g.Log().Errorf(r.Context(), "ChannelPaymentRedirectEntrance panic channelId: %s err:%s", r.GetUrl(), channelId, err)
		return
	}
	outchannel.GetPayChannelServiceProvider(channelIdInt).DoRemoteChannelRedirect(r)
}

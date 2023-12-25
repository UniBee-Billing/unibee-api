package webhooks

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/logic/payment/outchannel"
	"go-oversea-pay/internal/logic/payment/outchannel/util"
	"strconv"
)

func ChannelPaymentWebhookEntrance(r *ghttp.Request) {
	channelId := r.Get("channelId").String()
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		g.Log().Errorf(r.Context(), "ChannelPaymentWebhookEntrance panic channelId: %s err:%s", r.GetUrl(), channelId, err)
		return
	}
	payChannel := util.GetOverseaPayChannel(r.Context(), int64(channelIdInt))
	outchannel.GetPayChannelServiceProvider(r.Context(), int64(channelIdInt)).DoRemoteChannelWebhook(r, payChannel)
}

func ChannelPaymentRedirectEntrance(r *ghttp.Request) {
	channelId := r.Get("channelId").String()
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		g.Log().Errorf(r.Context(), "ChannelPaymentRedirectEntrance panic channelId: %s err:%s", r.GetUrl(), channelId, err)
		return
	}
	payChannel := util.GetOverseaPayChannel(r.Context(), int64(channelIdInt))
	outchannel.GetPayChannelServiceProvider(r.Context(), int64(channelIdInt)).DoRemoteChannelRedirect(r, payChannel)
}

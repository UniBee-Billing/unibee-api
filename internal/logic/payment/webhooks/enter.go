package webhooks

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/outchannel"
	entity "go-oversea-pay/internal/model/entity/oversea_pay"
	"strconv"
)

func GetPaymentWebhookEntranceUrl(pay *entity.OverseaPay) string {
	return fmt.Sprintf("%s/webhooks/%s/notifications?payId=%d", consts.GetConfigInstance().HostPath, pay.ChannelPayId, pay.Id)
}

func GetPaymentRedirectEntranceUrl(pay *entity.OverseaPay) string {
	return fmt.Sprintf("%s/redirect/%s/forward?payId=%d", consts.GetConfigInstance().HostPath, pay.ChannelPayId, pay.Id)
}

func ChannelPaymentWebhookEntrance(r *ghttp.Request) {
	channelId := r.Get("channelId").String()
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		g.Log().Errorf(r.Context(), "ChannelPaymentWebhookEntrance panic channelId: %s err:%s", r.GetUrl(), channelId, err)
		return
	}
	//outchannel := util.GetOverseaPayChannel(r.Context(), uint64(channelIdInt))
	outchannel.GetPayChannelServiceProvider(r.Context(), int64(channelIdInt)).DoRemoteChannelWebhook(r)
}

func ChannelPaymentRedirectEntrance(r *ghttp.Request) {
	channelId := r.Get("channelId").String()
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		g.Log().Errorf(r.Context(), "ChannelPaymentRedirectEntrance panic channelId: %s err:%s", r.GetUrl(), channelId, err)
		return
	}
	outchannel.GetPayChannelServiceProvider(r.Context(), int64(channelIdInt)).DoRemoteChannelRedirect(r)
}

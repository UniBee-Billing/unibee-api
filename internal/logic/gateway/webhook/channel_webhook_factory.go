package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func GetPayChannelWebhookServiceProvider(ctx context.Context, channelId int64) (channelService _interface.RemotePaymentChannelWebhookInterface) {
	proxy := &PayChannelWebhookProxy{}
	proxy.PaymentChannel = query.GetPayChannelById(ctx, channelId)
	utility.Assert(proxy.PaymentChannel != nil, fmt.Sprintf("channel not found %d", channelId))
	return proxy
}

func CheckAndSetupPayChannelWebhooks(ctx context.Context) {
	list := query.GetPayChannelsGroupByEnumKey(ctx)
	for _, paychannel := range list {
		err := GetPayChannelWebhookServiceProvider(ctx, int64(paychannel.Id)).DoRemoteChannelCheckAndSetupWebhook(ctx, paychannel)
		if err != nil {
			g.Log().Errorf(ctx, "CheckAndSetupPayChannelWebhooks GatewayName:%s Error:%s", paychannel.GatewayName, err)
		} else {
			g.Log().Infof(ctx, "CheckAndSetupPayChannelWebhooks GatewayName:%s Success", paychannel.GatewayName)
		}
	}
}

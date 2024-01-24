package channel

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func GetPayChannelServiceProvider(ctx context.Context, channelId int64) (channelService _interface.RemotePayChannelInterface) {
	proxy := &PayChannelProxy{}
	proxy.channel = query.GetPayChannelById(ctx, channelId)
	utility.Assert(proxy.channel != nil, fmt.Sprintf("channel not found %d", channelId))
	return proxy
}

func CheckAndSetupPayChannelWebhooks(ctx context.Context) {
	list := query.GetPayChannelsGroupByEnumKey(ctx)
	for _, paychannel := range list {
		err := GetPayChannelServiceProvider(ctx, int64(paychannel.Id)).DoRemoteChannelCheckAndSetupWebhook(ctx, paychannel)
		if err != nil {
			g.Log().Errorf(ctx, "CheckAndSetupPayChannelWebhooks channel:%s error:%s", paychannel.Channel, err)
		} else {
			g.Log().Infof(ctx, "CheckAndSetupPayChannelWebhooks channel:%s success", paychannel.Channel)
		}
	}
}

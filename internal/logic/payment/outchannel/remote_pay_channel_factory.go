package outchannel

import (
	"context"
	"go-oversea-pay/internal/query"
)

func GetPayChannelServiceProvider(ctx context.Context, channelId int64) (channelService RemotePayChannelInterface) {
	proxy := &PayChannelProxy{}
	proxy.channel = query.GetPaymentTypePayChannelById(ctx, channelId)
	return proxy
}

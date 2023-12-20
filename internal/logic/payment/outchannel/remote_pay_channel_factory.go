package outchannel

import (
	"context"
	"fmt"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func GetPayChannelServiceProvider(ctx context.Context, channelId int64) (channelService RemotePayChannelInterface) {
	proxy := &PayChannelProxy{}
	proxy.channel = query.GetPayChannelById(ctx, channelId)
	utility.Assert(proxy.channel != nil, fmt.Sprintf("channel not found %d", channelId))
	return proxy
}

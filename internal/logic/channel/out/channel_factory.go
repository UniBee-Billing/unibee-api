package out

import (
	"context"
	"fmt"
	"go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func GetPayChannelServiceProvider(ctx context.Context, channelId int64) (channelService _interface.RemotePayChannelInterface) {
	proxy := &PayChannelProxy{}
	proxy.PaymentChannel = query.GetPayChannelById(ctx, channelId)
	utility.Assert(proxy.PaymentChannel != nil, fmt.Sprintf("channel not found %d", channelId))
	return proxy
}

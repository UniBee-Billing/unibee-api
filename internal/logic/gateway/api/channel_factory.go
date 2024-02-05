package api

import (
	"context"
	"fmt"
	"go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func GetGatewayServiceProvider(ctx context.Context, gatewayId int64) (channelService _interface.RemotePayChannelInterface) {
	proxy := &PayChannelProxy{}
	proxy.PaymentChannel = query.GetGatewayById(ctx, gatewayId)
	utility.Assert(proxy.PaymentChannel != nil, fmt.Sprintf("gateway not found %d", gatewayId))
	return proxy
}

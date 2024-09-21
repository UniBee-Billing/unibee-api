package api

import (
	"context"
	"fmt"
	"unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"
)

func GetGatewayServiceProvider(ctx context.Context, gatewayId uint64) (one _interface.GatewayInterface) {
	proxy := &GatewayProxy{}
	proxy.Gateway = query.GetGatewayById(ctx, gatewayId)
	proxy.GatewayName = proxy.Gateway.GatewayName
	utility.Assert(proxy.Gateway != nil, fmt.Sprintf("gateway not found %d", gatewayId))
	return proxy
}

func GetGatewayWebhookServiceProviderByGatewayName(ctx context.Context, gatewayName string) (one _interface.GatewayInterface) {
	proxy := &GatewayProxy{}
	proxy.GatewayName = gatewayName
	return proxy
}

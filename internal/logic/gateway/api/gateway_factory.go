package api

import (
	"context"
	"fmt"
	"unibee-api/internal/interface"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func GetGatewayServiceProvider(ctx context.Context, gatewayId int64) (one _interface.GatewayInterface) {
	proxy := &GatewayProxy{}
	proxy.Gateway = query.GetGatewayById(ctx, gatewayId)
	utility.Assert(proxy.Gateway != nil, fmt.Sprintf("gateway not found %d", gatewayId))
	return proxy
}

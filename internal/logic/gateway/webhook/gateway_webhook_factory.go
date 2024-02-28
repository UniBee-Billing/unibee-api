package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/internal/interface"
	"unibee/internal/query"
	"unibee/utility"
)

func GetGatewayWebhookServiceProvider(ctx context.Context, gatewayId uint64) (one _interface.GatewayWebhookInterface) {
	proxy := &GatewayWebhookProxy{}
	proxy.Gateway = query.GetGatewayById(ctx, gatewayId)
	proxy.GatewayName = proxy.Gateway.GatewayName
	utility.Assert(proxy.Gateway != nil, fmt.Sprintf("gateway not found %d", gatewayId))
	return proxy
}

func GetGatewayWebhookServiceProviderByGatewayName(ctx context.Context, gatewayName string) (one _interface.GatewayWebhookInterface) {
	proxy := &GatewayWebhookProxy{}
	proxy.GatewayName = gatewayName
	return proxy
}

func CheckAndSetupGatewayWebhooks(ctx context.Context, gatewayId uint64) {
	gateway := query.GetGatewayById(ctx, gatewayId)
	utility.Assert(gateway != nil, fmt.Sprintf("gateway not found %d", gatewayId))
	err := GetGatewayWebhookServiceProvider(ctx, gateway.Id).GatewayCheckAndSetupWebhook(ctx, gateway)
	if err != nil {
		g.Log().Errorf(ctx, "CheckAndSetupGatewayWebhooks GatewayName:%s Error:%s", gateway.GatewayName, err)
	} else {
		g.Log().Infof(ctx, "CheckAndSetupGatewayWebhooks GatewayName:%s Success", gateway.GatewayName)
	}
	utility.AssertError(err, "CheckAndSetupGatewayWebhooks Error")
}

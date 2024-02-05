package webhook

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"
)

func GetGatewayWebhookServiceProvider(ctx context.Context, gatewayId int64) (one _interface.GatewayWebhookInterface) {
	proxy := &GatewayWebhookProxy{}
	proxy.PaymentChannel = query.GetGatewayById(ctx, gatewayId)
	utility.Assert(proxy.PaymentChannel != nil, fmt.Sprintf("gateway not found %d", gatewayId))
	return proxy
}

func CheckAndSetupGatewayWebhooks(ctx context.Context) {
	list := query.GetGatewaysGroupByEnumKey(ctx)
	for _, gateway := range list {
		err := GetGatewayWebhookServiceProvider(ctx, int64(gateway.Id)).GatewayCheckAndSetupWebhook(ctx, gateway)
		if err != nil {
			g.Log().Errorf(ctx, "CheckAndSetupGatewayWebhooks GatewayName:%s Error:%s", gateway.GatewayName, err)
		} else {
			g.Log().Infof(ctx, "CheckAndSetupGatewayWebhooks GatewayName:%s Success", gateway.GatewayName)
		}
	}
}

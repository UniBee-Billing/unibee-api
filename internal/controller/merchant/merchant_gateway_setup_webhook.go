package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	gateway2 "unibee/internal/logic/gateway"
	"unibee/internal/logic/gateway/api"
	gatewayWebhook "unibee/internal/logic/gateway/webhook"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) SetupWebhook(ctx context.Context, req *gateway.SetupWebhookReq) (res *gateway.SetupWebhookRes, err error) {
	one := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "merchant not match")
	gatewayInfo := api.GetGatewayServiceProvider(ctx, one.Id).GatewayInfo(ctx)
	utility.Assert(gatewayInfo != nil, "gateway not ready")
	gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)
	if len(gatewayInfo.GatewayWebhookIntegrationLink) > 0 && len(req.WebhookSecret) > 0 {
		err = query.UpdateGatewayWebhookSecret(ctx, one.Id, req.WebhookSecret)
		if err != nil {
			return nil, err
		}
	}

	return &gateway.SetupWebhookRes{GatewayWebhookUrl: gateway2.GetPaymentWebhookEntranceUrl(req.GatewayId)}, nil
}

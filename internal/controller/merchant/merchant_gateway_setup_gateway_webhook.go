package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	gatewayWebhook "unibee/internal/logic/gateway/webhook"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/gateway"
)

func (c *ControllerGateway) SetupGatewayWebhook(ctx context.Context, req *gateway.SetupGatewayWebhookReq) (res *gateway.SetupGatewayWebhookRes, err error) {
	one := query.GetGatewayById(ctx, req.GatewayId)
	utility.Assert(one != nil, "gateway not found")
	utility.Assert(one.MerchantId == _interface.GetMerchantId(ctx), "merchant not match")
	gatewayWebhook.CheckAndSetupGatewayWebhooks(ctx, one.Id)

	return &gateway.SetupGatewayWebhookRes{}, nil
}

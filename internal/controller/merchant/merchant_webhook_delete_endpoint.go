package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface/context"
	_webhook "unibee/internal/logic/webhook"
	"unibee/internal/query"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) DeleteEndpoint(ctx context.Context, req *webhook.DeleteEndpointReq) (res *webhook.DeleteEndpointRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = _webhook.DeleteMerchantWebhookEndpoint(ctx, _interface.GetMerchantId(ctx), req.EndpointId)
	if err != nil {
		return nil, err
	}
	return &webhook.DeleteEndpointRes{}, nil
}

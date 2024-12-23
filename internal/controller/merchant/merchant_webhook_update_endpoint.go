package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface/context"
	_webhook "unibee/internal/logic/webhook"
	"unibee/internal/query"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) UpdateEndpoint(ctx context.Context, req *webhook.UpdateEndpointReq) (res *webhook.UpdateEndpointRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = _webhook.UpdateMerchantWebhookEndpoint(ctx, _interface.GetMerchantId(ctx), req.EndpointId, req.Url, req.Events)
	if err != nil {
		return nil, err
	}
	return &webhook.UpdateEndpointRes{}, nil
}

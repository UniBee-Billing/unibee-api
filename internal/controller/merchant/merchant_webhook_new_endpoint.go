package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface"
	_webhook "unibee/internal/logic/webhook"
	"unibee/internal/query"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) NewEndpoint(ctx context.Context, req *webhook.NewEndpointReq) (res *webhook.NewEndpointRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	_, err = _webhook.NewMerchantWebhookEndpoint(ctx, _interface.GetMerchantId(ctx), req.Url, req.Events)
	if err != nil {
		return nil, err
	}
	return &webhook.NewEndpointRes{}, nil
}

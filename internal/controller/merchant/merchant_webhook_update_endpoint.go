package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	_webhook "unibee-api/internal/logic/webhook"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee-api/api/merchant/webhook"
)

func (c *ControllerWebhook) UpdateEndpoint(ctx context.Context, req *webhook.UpdateEndpointReq) (res *webhook.UpdateEndpointRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	one := query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = _webhook.UpdateMerchantWebhookEndpoint(ctx, _interface.GetMerchantId(ctx), req.EndpointId, req.Url, req.Events)
	if err != nil {
		return nil, err
	}
	return &webhook.UpdateEndpointRes{}, nil
}
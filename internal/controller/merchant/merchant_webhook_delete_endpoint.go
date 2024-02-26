package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	_webhook "unibee/internal/logic/webhook"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) DeleteEndpoint(ctx context.Context, req *webhook.DeleteEndpointReq) (res *webhook.DeleteEndpointRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	one := query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = _webhook.DeleteMerchantWebhookEndpoint(ctx, _interface.GetMerchantId(ctx), req.EndpointId)
	if err != nil {
		return nil, err
	}
	return &webhook.DeleteEndpointRes{}, nil
}

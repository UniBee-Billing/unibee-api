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

func (c *ControllerWebhook) UpdateEndpoint(ctx context.Context, req *webhook.UpdateEndpointReq) (res *webhook.UpdateEndpointRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
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

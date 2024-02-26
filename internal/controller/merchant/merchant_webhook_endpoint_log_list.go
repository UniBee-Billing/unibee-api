package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/webhook"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) EndpointLogList(ctx context.Context, req *webhook.EndpointLogListReq) (res *webhook.EndpointLogListRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	one := query.GetMerchantInfoById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	list := webhook2.MerchantWebhookEndpointLogList(ctx, &webhook2.EndpointLogListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		EndpointId: req.EndpointId,
		Page:       req.Page,
		Count:      req.Count,
	})
	return &webhook.EndpointLogListRes{EndpointLogList: list}, nil
}

package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	webhook2 "unibee-api/internal/logic/webhook"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"github.com/gogf/gf/v2/errors/gerror"

	"unibee-api/api/merchant/webhook"
)

func (c *ControllerWebhook) EndpointLogList(ctx context.Context, req *webhook.EndpointLogListReq) (res *webhook.EndpointLogListRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId > 0, "merchantUserId invalid")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.MerchantId == uint64(req.MerchantId), "merchantId not match")
	}
	one := query.GetMerchantInfoById(ctx, req.MerchantId)
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	list := webhook2.MerchantWebhookEndpointLogList(ctx, &webhook2.EndpointLogListInternalReq{
		MerchantId: req.MerchantId,
		EndpointId: req.EndpointId,
		Page:       req.Page,
		Count:      req.Count,
	})
	return &webhook.EndpointLogListRes{EndpointLogList: list}, nil
}

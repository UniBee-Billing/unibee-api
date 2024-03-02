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

func (c *ControllerWebhook) EndpointList(ctx context.Context, req *webhook.EndpointListReq) (res *webhook.EndpointListRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	list := webhook2.MerchantWebhookEndpointList(ctx, _interface.GetMerchantId(ctx))
	return &webhook.EndpointListRes{EndpointList: list}, nil
}

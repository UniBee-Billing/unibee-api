package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface/context"
	webhook2 "unibee/internal/logic/webhook"
	"unibee/internal/query"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) EndpointList(ctx context.Context, req *webhook.EndpointListReq) (res *webhook.EndpointListRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	list := webhook2.MerchantWebhookEndpointList(ctx, _interface.GetMerchantId(ctx))
	return &webhook.EndpointListRes{EndpointList: list}, nil
}

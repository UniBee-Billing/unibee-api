package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	_interface "unibee/internal/interface"
	webhook2 "unibee/internal/logic/webhook"
	"unibee/internal/query"

	"unibee/api/merchant/webhook"
)

func (c *ControllerWebhook) EndpointLogList(ctx context.Context, req *webhook.EndpointLogListReq) (res *webhook.EndpointLogListRes, err error) {
	one := query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	list, total := webhook2.MerchantWebhookEndpointLogList(ctx, &webhook2.EndpointLogListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		EndpointId: req.EndpointId,
		Page:       req.Page,
		Count:      req.Count,
	})
	return &webhook.EndpointLogListRes{EndpointLogList: list, Total: total}, nil
}

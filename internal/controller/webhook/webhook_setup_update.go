package webhook

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee-api/api/webhook/setup"
	_interface "unibee-api/internal/interface"
	_webhook "unibee-api/internal/logic/webhook"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

func (c *ControllerSetup) Update(ctx context.Context, req *setup.UpdateReq) (res *setup.UpdateRes, err error) {
	openApiConfig := _interface.BizCtx().Get(ctx).OpenApiConfig
	utility.Assert(openApiConfig != nil, "api config not found")
	utility.Assert(openApiConfig.MerchantId > 0, "api config not found")
	one := query.GetMerchantInfoById(ctx, openApiConfig.MerchantId)
	if one == nil {
		return nil, gerror.New("Merchant Check Error")
	}
	err = _webhook.SetupMerchantWebhook(ctx, openApiConfig.MerchantId, req.Url, req.Events)
	if err != nil {
		return nil, err
	}
	return &setup.UpdateRes{}, nil
}

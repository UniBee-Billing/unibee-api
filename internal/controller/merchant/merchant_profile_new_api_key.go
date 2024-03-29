package merchant

import (
	"context"
	"unibee/api/merchant/profile"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/merchant"
)

func (c *ControllerProfile) NewApiKey(ctx context.Context, req *profile.NewApiKeyReq) (res *profile.NewApiKeyRes, err error) {
	return &profile.NewApiKeyRes{ApiKey: merchant.NewOpenApiKey(ctx, _interface.GetMerchantId(ctx))}, nil
}

package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/multi_currencies"

	"unibee/api/merchant/profile"
)

func (c *ControllerProfile) SetupMultiCurrencies(ctx context.Context, req *profile.SetupMultiCurrenciesReq) (res *profile.SetupMultiCurrenciesRes, err error) {
	if req.MultiCurrencies != nil {
		multi_currencies.SetupMerchantMultiCurrenciesConfig(ctx, _interface.GetMerchantId(ctx), req.MultiCurrencies)
	}
	return &profile.SetupMultiCurrenciesRes{MultiCurrencies: multi_currencies.GetMerchantMultiCurrenciesConfig(ctx, _interface.GetMerchantId(ctx))}, nil
}

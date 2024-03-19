package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/profile"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/currency"
	"unibee/internal/query"
	"unibee/time"
)

func (c *ControllerProfile) Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error) {
	return &profile.GetRes{
		Merchant: bean.SimplifyMerchant(query.GetMerchantById(ctx, _interface.GetMerchantId(ctx))),
		Currency: currency.GetMerchantCurrencies(),
		Env:      config.GetConfigInstance().Env,
		IsProd:   config.GetConfigInstance().IsProd(),
		TimeZone: time.GetTimeZoneList(),
	}, nil
}

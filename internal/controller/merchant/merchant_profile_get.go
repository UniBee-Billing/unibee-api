package merchant

import (
	"context"
	"unibee/api/merchant/profile"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/currency"
	"unibee/internal/query"
	"unibee/time"
)

func (c *ControllerMerchantProfile) Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error) {
	return &profile.GetRes{
		Merchant: query.GetMerchantById(ctx, _interface.GetMerchantId(ctx)),
		Currency: currency.GetMerchantCurrencies(),
		Env:      consts.GetConfigInstance().Env,
		IsProd:   consts.GetConfigInstance().IsProd(),
		TimeZone: time.GetTimeZoneList(),
	}, nil
}

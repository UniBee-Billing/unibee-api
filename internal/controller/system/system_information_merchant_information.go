package system

import (
	"context"
	"unibee-api/api/system/information"
	"unibee-api/time"
)

func (c *ControllerInformation) MerchantInformation(ctx context.Context, req *information.MerchantInformationReq) (res *information.MerchantInformationRes, err error) {
	res = &information.MerchantInformationRes{}

	res.SupportTimeZone = time.GetTimeZoneList()

	var supportCurrencys []*information.SupportCurrency
	supportCurrencys = append(supportCurrencys, &information.SupportCurrency{
		Currency: "EUR",
		Symbol:   "€",
		Scale:    100,
	})
	supportCurrencys = append(supportCurrencys, &information.SupportCurrency{
		Currency: "USD",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrencys = append(supportCurrencys, &information.SupportCurrency{
		Currency: "JPY",
		Symbol:   "¥",
		Scale:    1,
	})
	res.SupportCurrency = supportCurrencys

	return res, nil
}

package system

import (
	"context"
	"unibee-api/api/system/information"
	"unibee-api/internal/consts"
	"unibee-api/internal/query"
	"unibee-api/time"
)

func (c *ControllerInformation) MerchantInformation(ctx context.Context, req *information.MerchantInformationReq) (res *information.MerchantInformationRes, err error) {
	res = &information.MerchantInformationRes{}

	res.SupportTimeZone = time.GetTimeZoneList()
	res.Env = consts.GetConfigInstance().Env
	res.IsProd = consts.GetConfigInstance().IsProd()

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
	res.MerchantId = 15621 // firstly only one
	res.MerchantInfo = query.GetMerchantInfoById(ctx, res.MerchantId)
	res.Gateways = query.GetListActiveOutGatewayRosByMerchantId(ctx, res.MerchantId)

	return res, nil
}

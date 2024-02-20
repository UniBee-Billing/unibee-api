package system

import (
	"context"
	"unibee-api/api/system/information"
	"unibee-api/internal/consts"
	"unibee-api/internal/logic/gateway"
	"unibee-api/internal/query"
	"unibee-api/time"
)

func (c *ControllerInformation) MerchantInformation(ctx context.Context, req *information.MerchantInformationReq) (res *information.MerchantInformationRes, err error) {
	res = &information.MerchantInformationRes{}

	res.SupportTimeZone = time.GetTimeZoneList()
	res.Env = consts.GetConfigInstance().Env
	res.IsProd = consts.GetConfigInstance().IsProd()

	var supportCurrency []*information.SupportCurrency
	supportCurrency = append(supportCurrency, &information.SupportCurrency{
		Currency: "EUR",
		Symbol:   "€",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &information.SupportCurrency{
		Currency: "USD",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &information.SupportCurrency{
		Currency: "JPY",
		Symbol:   "¥",
		Scale:    1,
	})
	res.SupportCurrency = supportCurrency
	res.MerchantId = 15621 // firstly only one
	res.MerchantInfo = query.GetMerchantInfoById(ctx, res.MerchantId)
	res.Gateway = gateway.GetListActiveOutGatewayRosByMerchantId(ctx, res.MerchantId)

	return res, nil
}

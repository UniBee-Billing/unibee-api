package system

import (
	"context"
	"unibee/api/system/information"
	"unibee/internal/consts"
	"unibee/time"
)

func (c *ControllerInformation) Get(ctx context.Context, req *information.GetReq) (res *information.GetRes, err error) {
	res = &information.GetRes{}

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

	return res, nil
}
package system

import (
	"context"
	"unibee/api/system/information"
	"unibee/internal/cmd/config"
	"unibee/internal/logic/currency"
	"unibee/time"
)

func (c *ControllerInformation) Get(ctx context.Context, req *information.GetReq) (res *information.GetRes, err error) {
	res = &information.GetRes{}

	res.SupportTimeZone = time.GetTimeZoneList()
	res.Env = config.GetConfigInstance().Env
	res.IsProd = config.GetConfigInstance().IsProd()
	res.SupportCurrency = currency.GetMerchantCurrencies()
	res.Mode = config.GetConfigInstance().Mode

	return res, nil
}

package system

import (
	"context"
	"unibee/api/system/information"
	"unibee/internal/consts"
	"unibee/internal/logic/currency"
	"unibee/time"
)

func (c *ControllerInformation) Get(ctx context.Context, req *information.GetReq) (res *information.GetRes, err error) {
	res = &information.GetRes{}

	res.SupportTimeZone = time.GetTimeZoneList()
	res.Env = consts.GetConfigInstance().Env
	res.IsProd = consts.GetConfigInstance().IsProd()
	res.SupportCurrency = currency.GetMerchantCurrencies()

	return res, nil
}

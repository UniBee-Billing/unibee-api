package merchant

import (
	"context"
	"strings"
	"unibee/api/bean"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/multi_currencies"
	"unibee/utility"

	"unibee/api/merchant/profile"
)

func (c *ControllerProfile) AmountMultiCurrenciesExchange(ctx context.Context, req *profile.AmountMultiCurrenciesExchangeReq) (res *profile.AmountMultiCurrenciesExchangeRes, err error) {
	req.Currency = strings.ToUpper(req.Currency)
	var multiCurrencies = make([]*bean.PlanMultiCurrency, 0)
	if len(req.Currency) > 0 && req.Amount > 0 {
		configs := multi_currencies.GetMerchantMultiCurrenciesConfig(ctx, _interface.GetMerchantId(ctx))
		for _, config := range configs {
			if strings.ToUpper(config.DefaultCurrency) == req.Currency {
				for _, multiCurrency := range config.MultiCurrencies {
					multiCurrencies = append(multiCurrencies, &bean.PlanMultiCurrency{
						Currency:     multiCurrency.Currency,
						AutoExchange: multiCurrency.AutoExchange,
						ExchangeRate: multiCurrency.ExchangeRate,
						Amount:       utility.ExchangeCurrencyConvert(req.Amount, req.Currency, multiCurrency.Currency, multiCurrency.ExchangeRate),
					})
				}

			}
		}
	}
	return &profile.AmountMultiCurrenciesExchangeRes{MultiCurrencies: multiCurrencies}, nil
}

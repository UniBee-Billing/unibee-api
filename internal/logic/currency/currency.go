package currency

import "unibee/internal/logic/gateway/ro"

func GetMerchantCurrencies() []*ro.Currency {
	var supportCurrency []*ro.Currency
	supportCurrency = append(supportCurrency, &ro.Currency{
		Currency: "EUR",
		Symbol:   "€",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &ro.Currency{
		Currency: "USD",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &ro.Currency{
		Currency: "JPY",
		Symbol:   "¥",
		Scale:    1,
	})
	return supportCurrency
}

package currency

import (
	"strings"
	"unibee/internal/logic/gateway/ro"
)

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

func GetMerchantCurrencyMap() map[string]*ro.Currency {
	var currencyMap = make(map[string]*ro.Currency)
	for _, currency := range GetMerchantCurrencies() {
		currencyMap[currency.Currency] = currency
	}
	return currencyMap
}

func IsFiatCurrencySupport(currency string) bool {
	//Fiat Currency Check
	if len(currency) == 0 {
		return false
	}
	return GetMerchantCurrencyMap()[strings.ToUpper(currency)] != nil
}

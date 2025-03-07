package currency

import (
	"strings"
	"unibee/api/bean"
)

func GetMerchantCurrencies() []*bean.Currency {
	var supportCurrency []*bean.Currency
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "EUR",
		Symbol:   "€",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "USD",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "RUB",
		Symbol:   "₽",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "JPY",
		Symbol:   "¥",
		Scale:    1,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "GBP",
		Symbol:   "£",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "CHF",
		Symbol:   "CHF",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "SEK",
		Symbol:   "kr",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "NOK",
		Symbol:   "kr",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "INR",
		Symbol:   "₹",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "CNY",
		Symbol:   "¥",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "SGD",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "HKD",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "THB",
		Symbol:   "฿",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "VND",
		Symbol:   "₫",
		Scale:    1,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "MXN",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "BRL",
		Symbol:   "R$",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "CAD",
		Symbol:   "$",
		Scale:    100,
	})
	supportCurrency = append(supportCurrency, &bean.Currency{
		Currency: "AUD",
		Symbol:   "$",
		Scale:    100,
	})
	return supportCurrency
}

func GetMerchantCurrencyMap() map[string]*bean.Currency {
	var currencyMap = make(map[string]*bean.Currency)
	for _, currency := range GetMerchantCurrencies() {
		currencyMap[currency.Currency] = currency
	}
	return currencyMap
}

func IsCurrencySupport(target string) bool {
	for _, currency := range GetMerchantCurrencies() {
		if strings.Compare(currency.Currency, strings.ToUpper(target)) == 0 {
			return true
		}
	}
	return false
}

func IsFiatCurrencySupport(currency string) bool {
	//Fiat Currency Check
	if len(currency) == 0 {
		return false
	}
	return GetMerchantCurrencyMap()[strings.ToUpper(currency)] != nil
}

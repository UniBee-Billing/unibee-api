package entity

import (
	"context"
	"fmt"
	"strings"
	dao "unibee/internal/dao/default"
	"unibee/utility"
)

type PlanMultiCurrency struct {
	Currency     string  `json:"currency" description:"target currency"`
	AutoExchange bool    `json:"autoExchange" description:"using https://app.exchangerate-api.com/ to update exchange rate if true, the exchange APIKey need setup first"`
	ExchangeRate float64 `json:"exchangeRate" description:"exchange rate, no setup required if AutoExchange is true"`
	Amount       int64   `json:"amount" description:"the amount of exchange rate"`
	Disable      bool    `json:"disable" description:"disable currency exchange"`
}

type MerchantCurrencyConfig struct {
	Currency     string  `json:"currency" description:"target currency"`
	AutoExchange bool    `json:"autoExchange" description:"using https://app.exchangerate-api.com/ to update exchange rate if true, the exchange APIKey need setup first"`
	ExchangeRate float64 `json:"exchangeRate"  description:"the exchange rate of gateway, no setup required if AutoExchange is true"`
}

type MerchantMultiCurrencyConfig struct {
	DefaultCurrency string                    `json:"defaultCurrency"`
	MultiCurrencies []*MerchantCurrencyConfig `json:"currencyConfigs"`
	LastUpdateTime  int64                     `json:"lastUpdateTime" description:"Last Update UTC Time"`
}

const MerchantMultiCurrenciesConfig = "MerchantMultiCurrenciesConfig"

func GetMerchantMultiCurrencyConfig(ctx context.Context, merchantId uint64) []*MerchantMultiCurrencyConfig {
	utility.Assert(merchantId > 0, "invalid merchantId")
	var one *MerchantConfig
	err := dao.MerchantConfig.Ctx(ctx).
		Where(dao.MerchantConfig.Columns().MerchantId, merchantId).
		Where(dao.MerchantConfig.Columns().ConfigKey, MerchantMultiCurrenciesConfig).
		Scan(&one)
	if err != nil {
		one = nil
	}
	var configs = make([]*MerchantMultiCurrencyConfig, 0)
	if one != nil {
		_ = utility.UnmarshalFromJsonString(one.ConfigValue, &configs)
	}
	return configs
}

func (p *Plan) CurrencyAmount(ctx context.Context, targetCurrency string) int64 {
	if targetCurrency == "" || strings.ToUpper(targetCurrency) == strings.ToUpper(p.Currency) {
		return p.Amount
	}

	var multiCurrencies = make([]*PlanMultiCurrency, 0)
	if len(p.GatewayProductDescription) > 0 {
		_ = utility.UnmarshalFromJsonString(p.GatewayProductDescription, &multiCurrencies)
	}
	for _, one := range multiCurrencies {
		if strings.ToUpper(one.Currency) == strings.ToUpper(targetCurrency) {
			if one.Disable {
				utility.Assert(false, fmt.Sprintf("Exchange: currency(%s) disabled for plan(%d)", targetCurrency, p.Id))
			}
		}
	}

	multiCurrencyConfigs := GetMerchantMultiCurrencyConfig(ctx, p.MerchantId)
	for _, multiCurrency := range multiCurrencyConfigs {
		if strings.ToUpper(multiCurrency.DefaultCurrency) == strings.ToUpper(p.Currency) {
			for _, one := range multiCurrency.MultiCurrencies {
				if strings.ToUpper(one.Currency) == strings.ToUpper(targetCurrency) {
					return utility.ExchangeCurrencyConvert(p.Amount, p.Currency, one.Currency, one.ExchangeRate)
				}
			}
		}
	}
	utility.Assert(false, fmt.Sprintf("CurrencyAmount: No currency(%s) config found for plan(%d)", targetCurrency, p.Id))
	return p.Amount
}

func (p *Plan) ExchangeAmountToCurrency(ctx context.Context, amount int64, targetCurrency string) int64 {
	if targetCurrency == "" || strings.ToUpper(targetCurrency) == strings.ToUpper(p.Currency) {
		return amount
	}

	var multiCurrencies = make([]*PlanMultiCurrency, 0)
	if len(p.GatewayProductDescription) > 0 {
		_ = utility.UnmarshalFromJsonString(p.GatewayProductDescription, &multiCurrencies)
	}
	for _, one := range multiCurrencies {
		if strings.ToUpper(one.Currency) == strings.ToUpper(targetCurrency) {
			if one.Disable {
				utility.Assert(false, fmt.Sprintf("Exchange: currency(%s) disabled for plan(%d)", targetCurrency, p.Id))
			}
		}
	}

	multiCurrencyConfigs := GetMerchantMultiCurrencyConfig(ctx, p.MerchantId)
	for _, multiCurrency := range multiCurrencyConfigs {
		if strings.ToUpper(multiCurrency.DefaultCurrency) == strings.ToUpper(p.Currency) {
			for _, one := range multiCurrency.MultiCurrencies {
				if strings.ToUpper(one.Currency) == strings.ToUpper(targetCurrency) {
					return utility.ExchangeCurrencyConvert(amount, p.Currency, one.Currency, one.ExchangeRate)
				}
			}
		}
	}
	utility.Assert(false, fmt.Sprintf("Exchange: No currency(%s) config found for plan(%d)", targetCurrency, p.Id))
	return amount
}

func ExchangeAmountToCurrency(ctx context.Context, merchantId uint64, amount int64, sourceCurrency string, targetCurrency string) int64 {
	if targetCurrency == "" || strings.ToUpper(targetCurrency) == strings.ToUpper(sourceCurrency) {
		return amount
	}
	multiCurrencyConfigs := GetMerchantMultiCurrencyConfig(ctx, merchantId)
	for _, multiCurrency := range multiCurrencyConfigs {
		if strings.ToUpper(multiCurrency.DefaultCurrency) == strings.ToUpper(sourceCurrency) {
			for _, one := range multiCurrency.MultiCurrencies {
				if strings.ToUpper(one.Currency) == strings.ToUpper(targetCurrency) {
					return utility.ExchangeCurrencyConvert(amount, sourceCurrency, one.Currency, one.ExchangeRate)
				}
			}
		}
	}
	utility.Assert(false, fmt.Sprintf("Exchange: No currency(%s) config found for merchant(%d)", targetCurrency, merchantId))
	return amount
}

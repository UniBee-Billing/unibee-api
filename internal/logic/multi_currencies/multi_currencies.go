package multi_currencies

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"strings"
	"unibee/api/bean"
	"unibee/internal/logic/currency"
	"unibee/internal/logic/merchant_config"
	"unibee/internal/logic/merchant_config/update"
	"unibee/internal/logic/multi_currencies/currency_exchange"
	"unibee/utility"
)

func SetupMerchantMultiCurrenciesConfig(ctx context.Context, merchantId uint64, configs []*bean.MerchantMultiCurrencyConfig) {
	for _, config := range configs {
		config.DefaultCurrency = strings.ToUpper(config.DefaultCurrency)
		utility.Assert(currency.IsCurrencySupport(config.DefaultCurrency), fmt.Sprintf("Default Currency is not supported:%s", config.DefaultCurrency))
		for _, exchange := range config.MultiCurrencies {
			exchange.Currency = strings.ToUpper(exchange.Currency)
			if exchange.AutoExchange {
				utility.Assert(currency.IsCurrencySupport(config.DefaultCurrency), fmt.Sprintf("Currency is not supported:%s", exchange.Currency))
				exchange.ExchangeRate = currency_exchange.GetMerchantExchangeCurrencyRate(ctx, merchantId, config.DefaultCurrency, exchange.Currency)
			} else {
				utility.Assert(exchange.ExchangeRate > 0, fmt.Sprintf("Invalid exchange rate, fromCurrency:%s toCurrency:%s", config.DefaultCurrency, exchange.Currency))
			}
		}
		config.LastUpdateTime = gtime.Now().Timestamp()
	}
	err := update.SetMerchantConfig(ctx, merchantId, currency_exchange.MerchantMultiCurrenciesConfig, utility.MarshalToJsonString(configs))
	utility.AssertError(err, "SetupMerchantMultiCurrenciesConfig")
}

func GetMerchantMultiCurrenciesConfig(ctx context.Context, merchantId uint64) []*bean.MerchantMultiCurrencyConfig {
	data := merchant_config.GetMerchantConfig(ctx, merchantId, currency_exchange.MerchantMultiCurrenciesConfig)
	var configs = make([]*bean.MerchantMultiCurrencyConfig, 0)
	if data != nil {
		_ = utility.UnmarshalFromJsonString(data.ConfigValue, &configs)
	}
	return configs
}

func UpdateMerchantMultiCurrenciesConfigExchangeRate(ctx context.Context, merchantId uint64) {
	configs := GetMerchantMultiCurrenciesConfig(ctx, merchantId)
	for _, config := range configs {
		utility.Assert(currency.IsCurrencySupport(config.DefaultCurrency), fmt.Sprintf("Default Currency is not supported:%s", config.DefaultCurrency))
		for _, exchange := range config.MultiCurrencies {
			if exchange.AutoExchange {
				utility.Assert(currency.IsCurrencySupport(config.DefaultCurrency), fmt.Sprintf("Currency is not supported:%s", exchange.Currency))
				exchange.ExchangeRate = currency_exchange.GetMerchantExchangeCurrencyRate(ctx, merchantId, config.DefaultCurrency, exchange.Currency)
			}
		}
		config.LastUpdateTime = gtime.Now().Timestamp()
	}
	err := update.SetMerchantConfig(ctx, merchantId, currency_exchange.MerchantMultiCurrenciesConfig, utility.MarshalToJsonString(configs))
	utility.AssertError(err, "SetupMerchantMultiCurrenciesConfig")
}

func GetMerchantMultiCurrenciesConfigMap(ctx context.Context, merchantId uint64) map[string][]*bean.MerchantCurrencyConfig {
	data := merchant_config.GetMerchantConfig(ctx, merchantId, currency_exchange.MerchantMultiCurrenciesConfig)
	var configs = make([]*bean.MerchantMultiCurrencyConfig, 0)
	if data != nil {
		_ = utility.UnmarshalFromJsonString(data.ConfigValue, &configs)
	}
	var configMap = make(map[string][]*bean.MerchantCurrencyConfig, 0)
	for _, config := range configs {
		configMap[config.DefaultCurrency] = config.MultiCurrencies
	}
	return configMap
}

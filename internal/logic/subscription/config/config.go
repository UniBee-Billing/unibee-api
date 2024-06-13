package config

import (
	"context"
	"strconv"
	"unibee/api/bean"
	"unibee/internal/logic/merchant_config"
)

const (
	DowngradeEffectImmediately         = "DowngradeEffectImmediately"
	UpdateProration                    = "UpgradeProration"
	IncompleteExpireTime               = "IncompleteExpireTime"
	InvoiceEmail                       = "InvoiceEmail"
	TryAutomaticPaymentBeforePeriodEnd = "TryAutomaticPaymentBeforePeriodEnd"
	FiatExchangeApiKey                 = "FiatExchangeApiKey"
	GatewayVATRule                     = "GatewayVATRule"
)

func GetMerchantSubscriptionConfig(ctx context.Context, merchantId uint64) (config *bean.SubscriptionConfig) {
	// default config
	config = &bean.SubscriptionConfig{
		DowngradeEffectImmediately:         false,
		UpgradeProration:                   true,
		IncompleteExpireTime:               24 * 60 * 60, // 24h expire after
		InvoiceEmail:                       true,
		TryAutomaticPaymentBeforePeriodEnd: 30 * 60, // 30 min before period
		GatewayVATRule:                     "",
	}
	downgradeEffectImmediatelyConfig := merchant_config.GetMerchantConfig(ctx, merchantId, DowngradeEffectImmediately)
	if downgradeEffectImmediatelyConfig != nil && downgradeEffectImmediatelyConfig.ConfigValue == "true" {
		config.DowngradeEffectImmediately = true
	}
	updateProrationConfig := merchant_config.GetMerchantConfig(ctx, merchantId, UpdateProration)
	if updateProrationConfig != nil && updateProrationConfig.ConfigValue == "false" {
		config.UpgradeProration = false
	}
	incompleteExpireTimeConfig := merchant_config.GetMerchantConfig(ctx, merchantId, IncompleteExpireTime)
	if incompleteExpireTimeConfig != nil && len(incompleteExpireTimeConfig.ConfigValue) > 0 {
		value, err := strconv.ParseInt(incompleteExpireTimeConfig.ConfigValue, 10, 64)
		if err == nil {
			config.IncompleteExpireTime = value
		}
	}
	invoiceEmailConfig := merchant_config.GetMerchantConfig(ctx, merchantId, InvoiceEmail)
	if invoiceEmailConfig != nil && invoiceEmailConfig.ConfigValue == "false" {
		config.InvoiceEmail = false
	}
	tryAutomaticPaymentBeforePeriodEnd := merchant_config.GetMerchantConfig(ctx, merchantId, TryAutomaticPaymentBeforePeriodEnd)
	if tryAutomaticPaymentBeforePeriodEnd != nil && len(incompleteExpireTimeConfig.ConfigValue) > 0 {
		value, err := strconv.ParseInt(incompleteExpireTimeConfig.ConfigValue, 10, 64)
		if err == nil {
			config.TryAutomaticPaymentBeforePeriodEnd = value
		}
	}
	vatREPLACERuleConfig := merchant_config.GetMerchantConfig(ctx, merchantId, GatewayVATRule)
	if vatREPLACERuleConfig != nil && len(vatREPLACERuleConfig.ConfigValue) > 0 {
		config.GatewayVATRule = vatREPLACERuleConfig.ConfigValue
	}
	return config
}

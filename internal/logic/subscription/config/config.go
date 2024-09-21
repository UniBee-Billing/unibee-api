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
	ShowZeroInvoice                    = "ShowZeroInvoice"
)

func GetMerchantSubscriptionConfig(ctx context.Context, merchantId uint64) (config *bean.SubscriptionConfig) {
	// default config
	config = &bean.SubscriptionConfig{
		DowngradeEffectImmediately:         false,
		UpgradeProration:                   true,
		IncompleteExpireTime:               24 * 60 * 60, // default 24h expire after
		InvoiceEmail:                       true,
		TryAutomaticPaymentBeforePeriodEnd: 2 * 60 * 60, // default 2 hours before period
		GatewayVATRule:                     "",
		ShowZeroInvoice:                    false, // default false
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
	gatewayVATRule := merchant_config.GetMerchantConfig(ctx, merchantId, GatewayVATRule)
	if gatewayVATRule != nil && len(gatewayVATRule.ConfigValue) > 0 {
		config.GatewayVATRule = gatewayVATRule.ConfigValue
	}
	showZeroInvoice := merchant_config.GetMerchantConfig(ctx, merchantId, ShowZeroInvoice)
	if showZeroInvoice != nil && showZeroInvoice.ConfigValue == "true" {
		config.ShowZeroInvoice = true
	}
	return config
}

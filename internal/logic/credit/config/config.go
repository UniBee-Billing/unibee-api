package config

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/internal/logic/credit/credit_query"
	currency2 "unibee/internal/logic/currency"
)

func CheckCreditConfig(ctx context.Context, merchantId uint64, creditType int, currency string) error {
	if merchantId <= 0 {
		return gerror.New("invalid merchantId")
	}
	if creditType != 1 && creditType != 2 {
		return gerror.New("invalid creditType")
	}
	if !currency2.IsCurrencySupport(currency) {
		return gerror.New("invalid currency")
	}
	one := credit_query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return gerror.New("credit config need setup")
	}
	return nil
}

func CheckCreditConfigRecurring(ctx context.Context, merchantId uint64, creditType int, currency string) bool {
	if merchantId <= 0 {
		return false
	}
	if creditType != 1 && creditType != 2 {
		return false
	}
	if !currency2.IsCurrencySupport(currency) {
		return false
	}
	one := credit_query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return false
	}
	if one.Recurring == 1 {
		return true
	}
	return false
}

func CheckCreditConfigPreviewDefaultUsed(ctx context.Context, merchantId uint64, creditType int, currency string) bool {
	if merchantId <= 0 {
		return false
	}
	if creditType != 1 && creditType != 2 {
		return false
	}
	if !currency2.IsCurrencySupport(currency) {
		return false
	}
	one := credit_query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return false
	}
	if one.PreviewDefaultUsed == 1 {
		return true
	}
	return false
}

func CheckCreditConfigDiscountCodeExclusive(ctx context.Context, merchantId uint64, creditType int, currency string) bool {
	if merchantId <= 0 {
		return false
	}
	if creditType != 1 && creditType != 2 {
		return false
	}
	if !currency2.IsCurrencySupport(currency) {
		return false
	}
	one := credit_query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return false
	}
	if one.DiscountCodeExclusive == 1 {
		return true
	}
	return false
}

func CheckCreditConfigRecharge(ctx context.Context, merchantId uint64, creditType int, currency string) error {
	if merchantId <= 0 {
		return gerror.New("invalid merchantId")
	}
	if creditType != 1 && creditType != 2 {
		return gerror.New("invalid creditType")
	}
	if !currency2.IsCurrencySupport(currency) {
		return gerror.New("invalid currency")
	}
	one := credit_query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return gerror.New("credit config need setup")
	}
	if one.RechargeEnable == 0 {
		return gerror.New("credit account recharge disable")
	}
	return nil
}

func CheckCreditConfigPayout(ctx context.Context, merchantId uint64, creditType int, currency string) error {
	if merchantId <= 0 {
		return gerror.New("invalid merchantId")
	}
	if creditType != 1 && creditType != 2 {
		return gerror.New("invalid creditType")
	}
	if !currency2.IsCurrencySupport(currency) {
		return gerror.New("invalid currency")
	}
	one := credit_query.GetCreditConfig(ctx, merchantId, creditType, currency)
	if one == nil {
		return gerror.New("credit config need setup")
	}
	if one.PayoutEnable == 0 {
		return gerror.New("credit account payout disable")
	}
	return nil
}

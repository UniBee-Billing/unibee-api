package credit_query

import (
	"context"
	"strings"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
)

func GetCreditConfigList(ctx context.Context, merchantId uint64, creditConfigType int) (list []*entity.CreditConfig) {
	if merchantId <= 0 {
		return nil
	}
	err := dao.CreditConfig.Ctx(ctx).
		Where(dao.CreditConfig.Columns().MerchantId, merchantId).
		Where(dao.CreditConfig.Columns().Type, creditConfigType).
		OrderAsc(dao.CreditConfig.Columns().Id).
		Scan(&list)
	if err != nil {
		return nil
	}
	return list
}

func GetUserAccountById(ctx context.Context, id uint64) (one *entity.UserAccount) {
	if id <= 0 {
		return nil
	}
	err := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetCreditAccountById(ctx context.Context, id uint64) (one *entity.CreditAccount) {
	if id <= 0 {
		return nil
	}
	err := dao.CreditAccount.Ctx(ctx).
		Where(dao.CreditAccount.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetCreditAccountByUserId(ctx context.Context, userId uint64, creditType int, currency string) (one *entity.CreditAccount) {
	if userId <= 0 {
		return nil
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if len(currency) == 0 {
		return nil
	}
	err := dao.CreditAccount.Ctx(ctx).
		Where(dao.CreditAccount.Columns().UserId, userId).
		Where(dao.CreditAccount.Columns().Type, creditType).
		Where(dao.CreditAccount.Columns().Currency, currency).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetCreditConfig(ctx context.Context, merchantId uint64, creditConfigType int, currency string) (one *entity.CreditConfig) {
	if merchantId <= 0 {
		return nil
	}
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if len(currency) == 0 {
		return nil
	}
	err := dao.CreditConfig.Ctx(ctx).
		Where(dao.CreditConfig.Columns().MerchantId, merchantId).
		Where(dao.CreditConfig.Columns().Type, creditConfigType).
		Where(dao.CreditConfig.Columns().Currency, currency).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetCreditPaymentByCreditPaymentId(ctx context.Context, creditPaymentId string) (one *entity.CreditPayment) {
	if len(creditPaymentId) <= 0 {
		return nil
	}
	err := dao.CreditPayment.Ctx(ctx).
		Where(dao.CreditPayment.Columns().CreditPaymentId, creditPaymentId).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetCreditRefundByExternalCreditRefundId(ctx context.Context, merchantId uint64, externalCreditRefundId string) (one *entity.CreditRefund) {
	if merchantId <= 0 {
		return nil
	}
	if len(externalCreditRefundId) <= 0 {
		return nil
	}
	err := dao.CreditRefund.Ctx(ctx).
		Where(dao.CreditRefund.Columns().MerchantId, merchantId).
		Where(dao.CreditRefund.Columns().ExternalCreditRefundId, externalCreditRefundId).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetCreditPaymentByExternalCreditPaymentId(ctx context.Context, merchantId uint64, externalCreditPaymentId string) (one *entity.CreditPayment) {
	if merchantId <= 0 {
		return nil
	}
	if len(externalCreditPaymentId) <= 0 {
		return nil
	}
	err := dao.CreditPayment.Ctx(ctx).
		Where(dao.CreditPayment.Columns().MerchantId, merchantId).
		Where(dao.CreditPayment.Columns().ExternalCreditPaymentId, externalCreditPaymentId).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

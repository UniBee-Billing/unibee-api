package query

import (
	"context"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

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
	one.Password = ""
	return one
}

func GetUserAccountByEmail(ctx context.Context, merchantId uint64, email string) (one *entity.UserAccount) {
	if len(email) == 0 {
		return nil
	}
	err := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().Email, email).
		Where(dao.UserAccount.Columns().MerchantId, merchantId).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetUserAccountByExternalUserId(ctx context.Context, merchantId uint64, externalUserId string) (one *entity.UserAccount) {
	if len(externalUserId) <= 0 {
		return nil
	}
	err := dao.UserAccount.Ctx(ctx).
		Where(dao.UserAccount.Columns().ExternalUserId, externalUserId).
		Where(dao.UserAccount.Columns().MerchantId, merchantId).
		Scan(&one)
	if err != nil {
		return nil
	}
	one.Password = ""
	return one
}

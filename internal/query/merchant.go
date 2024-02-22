package query

import (
	"context"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func GetMerchantInfoByApiKey(ctx context.Context, apiKey string) (one *entity.MerchantInfo) {
	if len(apiKey) <= 0 {
		return nil
	}
	err := dao.MerchantInfo.Ctx(ctx).
		Where(dao.MerchantInfo.Columns().ApiKey, apiKey).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantInfoById(ctx context.Context, id uint64) (one *entity.MerchantInfo) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantInfo.Ctx(ctx).
		Where(dao.MerchantInfo.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantInfoByHost(ctx context.Context, host string) (one *entity.MerchantInfo) {
	if len(host) <= 0 {
		return nil
	}
	err := dao.MerchantInfo.Ctx(ctx).
		Where(dao.MerchantInfo.Columns().Host, host).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetActiveMerchantInfoList(ctx context.Context) (list []*entity.MerchantInfo) {
	err := dao.MerchantInfo.Ctx(ctx).
		Where(dao.MerchantInfo.Columns().IsDeleted, 0).
		Scan(&list)
	if err != nil {
		return make([]*entity.MerchantInfo, 0)
	}
	return
}

func GetMerchantUserAccountById(ctx context.Context, id uint64) (one *entity.MerchantUserAccount) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantUserAccount.Ctx(ctx).
		Where(dao.MerchantUserAccount.Columns().Id, id).
		Scan(&one)
	if err != nil {
		return nil
	}
	one.Password = ""
	return one
}

func GetMerchantUserAccountByEmail(ctx context.Context, merchantId uint64, email string) (one *entity.MerchantUserAccount) {
	if len(email) == 0 {
		return nil
	}
	err := dao.MerchantUserAccount.Ctx(ctx).
		Where(dao.MerchantUserAccount.Columns().MerchantId, merchantId).
		Where(dao.MerchantUserAccount.Columns().Email, email).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

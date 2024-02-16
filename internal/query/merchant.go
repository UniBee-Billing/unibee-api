package query

import (
	"context"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func GetMerchantInfoById(ctx context.Context, id int64) (one *entity.MerchantInfo) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantInfo.Ctx(ctx).Where(entity.MerchantInfo{Id: id}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantAccountById(ctx context.Context, id uint64) (one *entity.MerchantUserAccount) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantUserAccount.Ctx(ctx).Where(entity.MerchantUserAccount{Id: id}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	one.Password = ""
	return one
}

func GetMerchantAccountByEmail(ctx context.Context, email string) (one *entity.MerchantUserAccount) {
	if len(email) == 0 {
		return nil
	}
	err := dao.MerchantUserAccount.Ctx(ctx).Where(entity.MerchantUserAccount{Email: email}).OmitEmpty().Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func GetMerchantCountryConfig(ctx context.Context, merchantId uint64, countryCode string) (one *entity.MerchantCountryConfig) {
	if merchantId <= 0 || len(countryCode) == 0 {
		return nil
	}
	err := dao.MerchantCountryConfig.Ctx(ctx).
		Where(dao.MerchantCountryConfig.Columns().MerchantId, merchantId).
		Where(dao.MerchantCountryConfig.Columns().CountryCode, countryCode).
		Scan(&one)
	if err != nil {
		one = nil
	}
	return
}

func GetMerchantCountryConfigName(ctx context.Context, merchantId uint64, countryCode string) string {
	one := GetMerchantCountryConfig(ctx, merchantId, countryCode)
	merchant := GetMerchantById(ctx, merchantId)
	utility.Assert(merchant != nil, "merchant not found")
	if one != nil {
		return one.Name
	}
	return merchant.Name
}

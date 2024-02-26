package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetOpenApiConfigByKey(ctx context.Context, key string) (res *entity.OpenApiConfig) {
	if len(key) == 0 {
		return nil
	}
	err := dao.OpenApiConfig.Ctx(ctx).Where(entity.OpenApiConfig{ApiKey: key}).OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res
}

func GetOneOpenApiConfigByMerchant(ctx context.Context, merchantId uint64) (res *entity.OpenApiConfig) {
	if merchantId <= 0 {
		return nil
	}
	err := dao.OpenApiConfig.Ctx(ctx).Where(entity.OpenApiConfig{MerchantId: merchantId}).OmitEmpty().Scan(&res)
	if err != nil {
		return nil
	}
	return res
}

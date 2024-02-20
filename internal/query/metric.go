package query

import (
	"context"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func GetMerchantMetric(ctx context.Context, id int64) (one *entity.MerchantMetric) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantMetric.Ctx(ctx).Where(entity.MerchantMetric{Id: uint64(id), IsDeleted: 0}).Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

func GetMerchantMetricByCode(ctx context.Context, code string) (one *entity.MerchantMetric) {
	if len(code) <= 0 {
		return nil
	}
	err := dao.MerchantMetric.Ctx(ctx).Where(entity.MerchantMetric{Code: code, IsDeleted: 0}).Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

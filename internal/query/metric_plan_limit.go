package query

import (
	"context"
	dao "unibee-api/internal/dao/oversea_pay"
	entity "unibee-api/internal/model/entity/oversea_pay"
)

func GetMerchantMetricPlanLimit(ctx context.Context, id int64) (one *entity.MerchantMetricPlanLimit) {
	if id <= 0 {
		return nil
	}
	err := dao.MerchantMetricPlanLimit.Ctx(ctx).
		Where(dao.MerchantMetricPlanLimit.Columns().Id, id).
		Where(dao.MerchantMetricPlanLimit.Columns().IsDeleted, 0).
		Scan(&one)
	if err != nil {
		return nil
	}
	return one
}

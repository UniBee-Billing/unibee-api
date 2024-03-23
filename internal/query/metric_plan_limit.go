package query

import (
	"context"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
)

func GetMerchantMetricPlanLimit(ctx context.Context, id uint64) (one *entity.MerchantMetricPlanLimit) {
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

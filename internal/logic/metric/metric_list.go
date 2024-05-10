package metric

import (
	"context"
	"unibee/api/bean"
	dao "unibee/internal/dao/oversea_pay"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/utility"
)

func MerchantMetricList(ctx context.Context, merchantId uint64) ([]*bean.MerchantMetricSimplify, int) {
	utility.Assert(merchantId > 0, "invalid merchantId")
	var list = make([]*bean.MerchantMetricSimplify, 0)
	if merchantId > 0 {
		var entities []*entity.MerchantMetric
		err := dao.MerchantMetric.Ctx(ctx).
			Where(dao.MerchantMetric.Columns().MerchantId, merchantId).
			Where(dao.MerchantMetric.Columns().IsDeleted, 0).
			Scan(&entities)
		if err == nil && len(entities) > 0 {
			for _, one := range entities {
				list = append(list, bean.SimplifyMerchantMetric(one))
			}
		}
	}
	return list, len(list)
}

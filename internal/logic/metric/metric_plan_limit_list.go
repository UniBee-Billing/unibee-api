package metric

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/utility"
)

func MerchantMetricPlanLimitCachedList(ctx context.Context, merchantId uint64, planId uint64, reloadCache bool) []*bean.MerchantMetricPlanLimit {
	utility.Assert(merchantId > 0, "invalid merchantId")
	utility.Assert(planId > 0, "invalid planId")
	var list = make([]*bean.MerchantMetricPlanLimit, 0)
	cacheKey := fmt.Sprintf("%s%d%d", MerchantMetricPlanLimitCacheKeyPrefix, merchantId, planId)
	if !reloadCache {
		get, err := g.Redis().Get(ctx, cacheKey)
		if err == nil && !get.IsNil() && !get.IsEmpty() {
			value := get.String()
			_ = utility.UnmarshalFromJsonString(value, &list)
			if len(list) > 0 {
				return list
			}
		}
	}
	if merchantId > 0 {
		var entities []*entity.MerchantMetricPlanLimit
		err := dao.MerchantMetricPlanLimit.Ctx(ctx).
			Where(dao.MerchantMetricPlanLimit.Columns().MerchantId, merchantId).
			Where(dao.MerchantMetricPlanLimit.Columns().PlanId, planId).
			Where(dao.MerchantMetricPlanLimit.Columns().IsDeleted, 0).
			Scan(&entities)
		if err == nil && len(entities) > 0 {
			for _, one := range entities {
				list = append(list, &bean.MerchantMetricPlanLimit{
					Id:          one.Id,
					MerchantId:  one.MerchantId,
					MetricId:    one.MetricId,
					Metric:      GetMerchantMetricSimplify(ctx, one.MetricId),
					PlanId:      one.PlanId,
					MetricLimit: one.MetricLimit,
					UpdateTime:  one.GmtModify.Timestamp(),
					CreateTime:  one.CreateTime,
				})
			}
		}
	}
	if len(list) > 0 {
		_, _ = g.Redis().Set(ctx, cacheKey, utility.MarshalToJsonString(list))
		_, _ = g.Redis().Expire(ctx, cacheKey, MerchantMetricPlanLimitCacheExpire) // one day cache expire time
	}
	return list
}

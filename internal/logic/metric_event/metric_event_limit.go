package metric_event

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/metric"
	"unibee-api/internal/logic/subscription/user_sub_plan"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/utility"
)

type MetricLimitVo struct {
	MerchantId          int64
	UserId              int64
	MetricId            int64
	Code                string `json:"code"                description:"code"`                                                                        // code
	MetricName          string `json:"metricName"          description:"metric name"`                                                                 // metric name
	Type                int    `json:"type"                description:"1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)"` // 1-limit_metered，2-charge_metered(come later),3-charge_recurring(come later)
	AggregationType     int    `json:"aggregationType"     description:"0-count，1-count unique, 2-latest, 3-max, 4-sum"`                              // 0-count，1-count unique, 2-latest, 3-max, 4-sum
	AggregationProperty string `json:"aggregationProperty" description:"aggregation property"`
	TotalLimit          uint64
	PlanLimits          []*ro.MerchantMetricPlanLimitVo // ?
}

func checkMetricLimitReached(ctx context.Context, merchantId int64, user *entity.UserAccount, met *entity.MerchantMetric, append uint64) (uint64, uint64, bool) {
	limitMap := GetUserMetricTotalLimits(ctx, merchantId, int64(user.Id))
	if metricLimit, ok := limitMap[int64(met.Id)]; ok {
		useValue := GetUserMetricLimitCachedUseValue(ctx, merchantId, int64(user.Id), met, false)
		if met.AggregationType == metric.MetricAggregationTypeLatest || met.AggregationType == metric.MetricAggregationTypeMax {
			return useValue, metricLimit.TotalLimit, append <= metricLimit.TotalLimit
		} else {
			return useValue, metricLimit.TotalLimit, useValue+append <= metricLimit.TotalLimit
		}
	} else {
		// no limit found, reject
		return 0, 0, false
	}
}

func GetUserMetricTotalLimits(ctx context.Context, merchantId int64, userId int64) map[int64]*MetricLimitVo {
	var limitMap = make(map[int64]*MetricLimitVo)
	userSubPlans := user_sub_plan.UserSubPlanCachedList(ctx, merchantId, userId, false)
	if len(userSubPlans) > 0 {
		for _, subPlan := range userSubPlans {
			list := metric.MerchantMetricPlanLimitCachedList(ctx, merchantId, subPlan.PlanId, false)
			for _, planLimit := range list {
				if _, ok := limitMap[planLimit.MetricId]; ok {
					limitMap[planLimit.MetricId].TotalLimit = limitMap[planLimit.MetricId].TotalLimit + planLimit.MetricLimit
					limitMap[planLimit.MetricId].PlanLimits = append(limitMap[planLimit.MetricId].PlanLimits, planLimit)
				} else {
					limitMap[planLimit.MetricId] = &MetricLimitVo{
						MerchantId:          merchantId,
						UserId:              userId,
						MetricId:            planLimit.MetricId,
						Code:                planLimit.Metric.Code,
						MetricName:          planLimit.Metric.MetricName,
						Type:                planLimit.Metric.Type,
						AggregationType:     planLimit.Metric.AggregationType,
						AggregationProperty: planLimit.Metric.AggregationProperty,
						TotalLimit:          uint64(planLimit.MetricLimit),
						PlanLimits:          []*ro.MerchantMetricPlanLimitVo{planLimit},
					}
				}
			}
		}
	}
	return limitMap
}

const (
	UserMetricCacheKeyPrefix = "UserMetricCacheKeyPrefix_"
	UserMetricCacheKeyExpire = 15 * 24 * 60 * 60 // 15 days cache expire
)

func GetUserMetricLimitCachedUseValue(ctx context.Context, merchantId int64, userId int64, met *entity.MerchantMetric, reloadCache bool) uint64 {
	cacheKey := fmt.Sprintf("%s_%d_%d_%d", UserMetricCacheKeyPrefix, merchantId, userId, met.Id)
	if !reloadCache {
		get, err := g.Redis().Get(ctx, cacheKey)
		if err == nil && !get.IsNil() && !get.IsEmpty() && (get.IsUint() || get.IsInt()) {
			return get.Uint64()
		}
	}
	var useValue uint64 = 0

	if merchantId > 0 {
		// count useValue from database
		if met.AggregationType == metric.MetricAggregationTypeLatest {
			useValue = 0 // type of this not need to compute from db
		} else if met.AggregationType == metric.MetricAggregationTypeMax {
			useValueFloat, err := dao.MerchantMetricEvent.Ctx(ctx).
				Where(entity.MerchantMetricEvent{MerchantId: merchantId}).
				Where(entity.MerchantMetricEvent{UserId: userId}).
				Where(entity.MerchantMetricEvent{MetricId: int64(met.Id)}).
				Max(dao.MerchantMetricEvent.Columns().AggregationPropertyInt)
			utility.AssertError(err, "server err")
			useValue = uint64(useValueFloat)
		} else {
			useValueFloat, err := dao.MerchantMetricEvent.Ctx(ctx).
				Where(entity.MerchantMetricEvent{MerchantId: merchantId}).
				Where(entity.MerchantMetricEvent{UserId: userId}).
				Where(entity.MerchantMetricEvent{MetricId: int64(met.Id)}).
				Sum(dao.MerchantMetricEvent.Columns().AggregationPropertyInt)
			utility.AssertError(err, "server err")
			useValue = uint64(useValueFloat)
		}
	}

	_, _ = g.Redis().Set(ctx, cacheKey, useValue)
	_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)

	return useValue
}

func appendMetricLimitCachedUseValue(ctx context.Context, merchantId int64, user *entity.UserAccount, met *entity.MerchantMetric, append uint64) {
	cacheKey := fmt.Sprintf("%s_%d_%d_%d", UserMetricCacheKeyPrefix, merchantId, user.Id, met.Id)
	get, err := g.Redis().Get(ctx, cacheKey)
	if err == nil && !get.IsNil() && !get.IsEmpty() && (get.IsUint() || get.IsInt()) {
		newValue := get.Uint64() + append
		if met.AggregationType == metric.MetricAggregationTypeLatest {
			newValue = append
		} else if met.AggregationType == metric.MetricAggregationTypeMax {
			newValue = utility.MaxUInt64(get.Uint64(), append)
		}
		_, _ = g.Redis().Set(ctx, cacheKey, newValue)
		_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)
	} else {
		_, _ = g.Redis().Set(ctx, cacheKey, append)
		_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)
	}
}

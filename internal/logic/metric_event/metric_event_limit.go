package metric_event

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	metric2 "unibee/api/merchant/metric"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/metric"
	addon2 "unibee/internal/logic/subscription/addon"
	"unibee/internal/logic/subscription/user_sub_plan"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func GetUserMetricStat(ctx context.Context, merchantId uint64, user *entity.UserAccount, productId int64) *metric2.UserMetric {
	sub := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, user.Id, user.MerchantId, productId)
	if sub == nil {
		sub = query.GetLatestSubscriptionByUserId(ctx, user.Id, user.MerchantId, productId)
	}
	return GetUserSubscriptionMetricStat(ctx, merchantId, user, sub)
}

func GetUserSubscriptionMetricStat(ctx context.Context, merchantId uint64, user *entity.UserAccount, one *entity.Subscription) *metric2.UserMetric {
	var list = make([]*bean.UserMerchantMetricStat, 0)
	if one != nil {
		limitMap := GetUserMetricTotalLimits(ctx, merchantId, user.Id, one)
		for _, metricLimit := range limitMap {
			met := query.GetMerchantMetric(ctx, metricLimit.MetricId)
			if met != nil {
				list = append(list, &bean.UserMerchantMetricStat{
					MetricLimit:     metricLimit,
					CurrentUseValue: GetUserMetricLimitCachedUseValue(ctx, merchantId, user.Id, met, one, false),
				})
			}
		}
		return &metric2.UserMetric{
			IsPaid:                  one.Status == consts.SubStatusActive || one.Status == consts.SubStatusIncomplete,
			Product:                 bean.SimplifyProduct(query.GetProductById(ctx, uint64(query.GetPlanById(ctx, one.PlanId).ProductId), merchantId)),
			User:                    bean.SimplifyUserAccount(user),
			Subscription:            bean.SimplifySubscription(ctx, one),
			Plan:                    bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:                  addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UserMerchantMetricStats: list,
		}
	} else {
		return &metric2.UserMetric{
			IsPaid:                  false,
			User:                    bean.SimplifyUserAccount(user),
			Product:                 nil,
			Subscription:            nil,
			Plan:                    nil,
			Addons:                  nil,
			UserMerchantMetricStats: list,
		}
	}
}

func checkMetricLimitReached(ctx context.Context, merchantId uint64, user *entity.UserAccount, sub *entity.Subscription, met *entity.MerchantMetric, append uint64) (uint64, uint64, bool) {
	limitMap := GetUserMetricTotalLimits(ctx, merchantId, user.Id, sub)
	if metricLimit, ok := limitMap[met.Id]; ok {
		useValue := GetUserMetricLimitCachedUseValue(ctx, merchantId, user.Id, met, sub, false)
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

func GetUserMetricTotalLimits(ctx context.Context, merchantId uint64, userId uint64, sub *entity.Subscription) map[uint64]*bean.PlanMetricLimitDetail {
	var limitMap = make(map[uint64]*bean.PlanMetricLimitDetail)
	userSubPlans := user_sub_plan.UserSubPlanCachedListForMetric(ctx, merchantId, userId, sub, false)
	if len(userSubPlans) > 0 {
		g.Log().Infof(ctx, "GetUserMetricTotalLimits userId:%d subPlanId:%d userSubPlans:%s", userId, sub.PlanId, utility.MarshalToJsonString(userSubPlans))
		for _, subPlan := range userSubPlans {
			list := metric.MerchantMetricPlanLimitCachedList(ctx, merchantId, subPlan.PlanId, false)
			for _, planLimit := range list {
				if _, ok := limitMap[planLimit.MetricId]; ok {
					limitMap[planLimit.MetricId].TotalLimit = limitMap[planLimit.MetricId].TotalLimit + planLimit.MetricLimit
					limitMap[planLimit.MetricId].PlanLimits = append(limitMap[planLimit.MetricId].PlanLimits, planLimit)
				} else {
					limitMap[planLimit.MetricId] = &bean.PlanMetricLimitDetail{
						MerchantId:          merchantId,
						UserId:              userId,
						MetricId:            planLimit.MetricId,
						Code:                planLimit.Metric.Code,
						MetricName:          planLimit.Metric.MetricName,
						Type:                planLimit.Metric.Type,
						AggregationType:     planLimit.Metric.AggregationType,
						AggregationProperty: planLimit.Metric.AggregationProperty,
						TotalLimit:          planLimit.MetricLimit,
						PlanLimits:          []*bean.MerchantMetricPlanLimit{planLimit},
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

func GetUserMetricLimitCachedUseValue(ctx context.Context, merchantId uint64, userId uint64, met *entity.MerchantMetric, sub *entity.Subscription, reloadCache bool) uint64 {
	cacheKey := metricUserCacheKey(merchantId, userId, met, sub)
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
			var latestOne *entity.MerchantMetricEvent
			q := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				OrderDesc(dao.MerchantMetricEvent.Columns().GmtCreate)
			if met.Type != metric.MetricTypeChargeRecurring {
				q = q.Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
					Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd)
			}
			err := q.Scan(&latestOne)
			utility.AssertError(err, "Server Error")
			if latestOne != nil {
				useValue = latestOne.AggregationPropertyInt
			}
		} else if met.AggregationType == metric.MetricAggregationTypeMax {
			q := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0)
			if met.Type != metric.MetricTypeChargeRecurring {
				q = q.Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
					Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd)
			}
			useValueFloat, err := q.Max(dao.MerchantMetricEvent.Columns().AggregationPropertyInt)
			utility.AssertError(err, "Server Error")
			useValue = uint64(useValueFloat)
		} else {
			q := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0)
			if met.Type != metric.MetricTypeChargeRecurring {
				q = q.Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
					Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd)
			}
			useValueFloat, err := q.Sum(dao.MerchantMetricEvent.Columns().AggregationPropertyInt)
			utility.AssertError(err, "Server Error")
			useValue = uint64(useValueFloat)
		}
	}

	_, _ = g.Redis().Set(ctx, cacheKey, useValue)
	_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)

	return useValue
}

func appendMetricLimitCachedUseValue(ctx context.Context, merchantId uint64, user *entity.UserAccount, met *entity.MerchantMetric, sub *entity.Subscription, append uint64) uint64 {
	cacheKey := metricUserCacheKey(merchantId, user.Id, met, sub)
	get, err := g.Redis().Get(ctx, cacheKey)
	if err == nil && !get.IsNil() && !get.IsEmpty() {
		newValue := get.Uint64() + append
		if met.AggregationType == metric.MetricAggregationTypeLatest {
			newValue = append
		} else if met.AggregationType == metric.MetricAggregationTypeMax {
			newValue = utility.MaxUInt64(get.Uint64(), append)
		}
		_, _ = g.Redis().Set(ctx, cacheKey, newValue)
		_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)
		return newValue
	} else {
		_, _ = g.Redis().Set(ctx, cacheKey, append)
		_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)
		return append
	}
}

func metricUserCacheKey(merchantId uint64, userId uint64, met *entity.MerchantMetric, sub *entity.Subscription) string {
	cacheKey := fmt.Sprintf("%s_%d_%d_%d_%s_%d", UserMetricCacheKeyPrefix, merchantId, userId, met.Id, sub.SubscriptionId, sub.CurrentPeriodStart)
	return cacheKey
}

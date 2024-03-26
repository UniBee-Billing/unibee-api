package metric_event

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean"
	metric2 "unibee/api/merchant/metric"
	"unibee/internal/consts"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/metric"
	addon2 "unibee/internal/logic/subscription/addon"
	"unibee/internal/logic/subscription/user_sub_plan"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

func GetUserMetricStat(ctx context.Context, merchantId uint64, user *entity.UserAccount) *metric2.UserMetric {
	sub := query.GetSubscriptionBySubscriptionId(ctx, user.SubscriptionId)
	return GetUserSubscriptionMetricStat(ctx, merchantId, user, sub)
}

func GetUserSubscriptionMetricStat(ctx context.Context, merchantId uint64, user *entity.UserAccount, one *entity.Subscription) *metric2.UserMetric {
	var list = make([]*bean.UserMerchantMetricStat, 0)
	if user != nil {
		user.Password = ""
	}
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
			User:                    bean.SimplifyUserAccount(user),
			Subscription:            bean.SimplifySubscription(one),
			Plan:                    bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:                  addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			UserMerchantMetricStats: list,
		}
	} else {
		return &metric2.UserMetric{
			IsPaid:                  false,
			User:                    bean.SimplifyUserAccount(user),
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
		for _, subPlan := range userSubPlans {
			list := metric.MerchantMetricPlanLimitCachedList(ctx, merchantId, subPlan.PlanId, true)
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

func ReloadUserMetricLimitCacheBackground(ctx context.Context, merchantId uint64, userId uint64, metricId uint64) {
	go func() {
		ctx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				g.Log().Errorf(ctx, "ReloadUserSubPlanCacheListBackground panic error:%s", err.Error())
				return
			}
		}()
		met := query.GetMerchantMetric(ctx, metricId)
		if met != nil {
			sub := query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, userId, merchantId)
			if sub != nil {
				GetUserMetricLimitCachedUseValue(ctx, merchantId, uint64(userId), met, sub, true)
			}
		}
	}()
}

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
			err := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0).
				OrderDesc(dao.MerchantMetricEvent.Columns().GmtCreate).
				Scan(&latestOne)
			utility.AssertError(err, "Server Error")
			if latestOne != nil {
				useValue = latestOne.AggregationPropertyInt
			}
		} else if met.AggregationType == metric.MetricAggregationTypeMax {
			useValueFloat, err := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0).
				Max(dao.MerchantMetricEvent.Columns().AggregationPropertyInt)
			utility.AssertError(err, "Server Error")
			useValue = uint64(useValueFloat)
		} else {
			useValueFloat, err := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0).
				Sum(dao.MerchantMetricEvent.Columns().AggregationPropertyInt)
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

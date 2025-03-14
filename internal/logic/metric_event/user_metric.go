package metric_event

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/log"
	dao "unibee/internal/dao/default"
	"unibee/internal/logic/metric"
	"unibee/internal/logic/metric_event/event_charge"
	addon2 "unibee/internal/logic/subscription/addon"
	"unibee/internal/logic/subscription/user_sub_plan"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

func UpdateMetricEventForInvoicePaid(one *entity.Invoice) {
	if one == nil {
		return
	}
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
				log.PrintPanic(ctx, err)
				return
			}
		}()
		invoiceDetail := detail.ConvertInvoiceToDetail(ctx, one)
		var list = make([]*bean.UserMetricChargeInvoiceItem, 0)
		if invoiceDetail.UserMetricChargeForInvoice != nil && len(invoiceDetail.UserMetricChargeForInvoice.MeteredChargeStats) > 0 {
			list = append(list, invoiceDetail.UserMetricChargeForInvoice.MeteredChargeStats...)
		}
		if invoiceDetail.UserMetricChargeForInvoice != nil && len(invoiceDetail.UserMetricChargeForInvoice.RecurringChargeStats) > 0 {
			list = append(list, invoiceDetail.UserMetricChargeForInvoice.RecurringChargeStats...)
		}
		for _, item := range list {
			if item != nil && item.MetricId > 0 && item.MaxEventId > 0 {
				_, err = dao.MerchantMetricEvent.Ctx(ctx).Data(g.Map{
					dao.MerchantMetricEvent.Columns().ChargeStatus:    1,
					dao.MerchantMetricEvent.Columns().ChargeInvoiceId: one.InvoiceId,
					dao.MerchantMetricEvent.Columns().GmtModify:       gtime.Now(),
				}).Where(dao.MerchantMetricEvent.Columns().ChargeStatus, 0).
					Where(dao.MerchantMetricEvent.Columns().MerchantId, one.MerchantId).
					Where(dao.MerchantMetricEvent.Columns().UserId, one.UserId).
					Where(dao.MerchantMetricEvent.Columns().MetricId, item.MetricId).
					WhereLTE(dao.MerchantMetricEvent.Columns().Id, item.MaxEventId).
					WhereGTE(dao.MerchantMetricEvent.Columns().Id, item.MinEventId).
					Update()
				if err != nil {
					g.Log().Errorf(ctx, "Update MetricEvent for invoice paid, invoiceId:%s metricId:%d maxEventId:%d err:%s", one.InvoiceId, item.MetricId, item.MaxEventId, err.Error())
				}
			}
		}
	}()
}

func GetUserMetricStatForAutoChargeInvoice(ctx context.Context, merchantId uint64, user *entity.UserAccount, sub *entity.Subscription, reloadCache bool) (entity *bean.UserMetricChargeInvoiceItemEntity) {
	entity = &bean.UserMetricChargeInvoiceItemEntity{
		MeteredChargeStats:   make([]*bean.UserMetricChargeInvoiceItem, 0),
		RecurringChargeStats: make([]*bean.UserMetricChargeInvoiceItem, 0),
	}
	defer func() {
		if exception := recover(); exception != nil {
			return
		}
	}()
	userMetric := GetUserSubscriptionMetricStat(ctx, merchantId, user, sub, reloadCache)
	meteredChargeStats := make([]*bean.UserMetricChargeInvoiceItem, 0)
	recurringChargeStats := make([]*bean.UserMetricChargeInvoiceItem, 0)
	for _, v := range userMetric.MeteredChargeStats {
		meteredChargeStats = append(meteredChargeStats, &bean.UserMetricChargeInvoiceItem{
			MetricId:          v.MetricId,
			CurrentUsedValue:  v.CurrentUsedValue,
			MaxEventId:        v.MaxEventId,
			MinEventId:        v.MinEventId,
			ChargePricing:     v.ChargePricing,
			TotalChargeAmount: v.TotalChargeAmount,
			Name:              v.Metric.MetricName,
			Description:       v.Metric.MetricDescription,
		})
	}
	for _, v := range userMetric.RecurringChargeStats {
		recurringChargeStats = append(recurringChargeStats, &bean.UserMetricChargeInvoiceItem{
			MetricId:          v.MetricId,
			CurrentUsedValue:  v.CurrentUsedValue,
			MaxEventId:        v.MaxEventId,
			MinEventId:        v.MinEventId,
			ChargePricing:     v.ChargePricing,
			TotalChargeAmount: v.TotalChargeAmount,
			Name:              v.Metric.MetricName,
			Description:       v.Metric.MetricDescription,
		})
	}
	return &bean.UserMetricChargeInvoiceItemEntity{
		MeteredChargeStats:   meteredChargeStats,
		RecurringChargeStats: recurringChargeStats,
	}
}

func GetUserMetricStat(ctx context.Context, merchantId uint64, user *entity.UserAccount, productId int64, reloadCache bool) *detail.UserMetric {
	sub := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, user.Id, user.MerchantId, productId)
	if sub == nil {
		sub = query.GetLatestSubscriptionByUserId(ctx, user.Id, user.MerchantId, productId)
	}
	return GetUserSubscriptionMetricStat(ctx, merchantId, user, sub, reloadCache)
}

func GetUserSubscriptionMetricStat(ctx context.Context, merchantId uint64, user *entity.UserAccount, one *entity.Subscription, reloadCache bool) *detail.UserMetric {
	var list = make([]*detail.UserMerchantMetricLimitStat, 0)
	var meteredChargeStats = make([]*detail.UserMerchantMetricChargeStat, 0)
	var recurringChargeStats = make([]*detail.UserMerchantMetricChargeStat, 0)
	if one != nil {
		limitMap := GetUserMetricTotalLimits(ctx, merchantId, user.Id, one)
		for _, metricLimit := range limitMap {
			met := query.GetMerchantMetric(ctx, metricLimit.MetricId)
			if met != nil {
				list = append(list, &detail.UserMerchantMetricLimitStat{
					MetricLimit:      metricLimit,
					CurrentUsedValue: GetUserMetricCachedUseValue(ctx, merchantId, user.Id, met, one, reloadCache).UsedValue,
				})
			}
		}
		planMetricBindingEntity := bean.ConvertMetricPlanBindingEntityFromPlan(query.GetPlanById(ctx, one.PlanId))
		for _, metricMeteredCharge := range planMetricBindingEntity.MetricMeteredCharge {
			met := query.GetMerchantMetric(ctx, metricMeteredCharge.MetricId)
			metricUserUsedValue := GetUserMetricCachedUseValue(ctx, merchantId, user.Id, met, one, reloadCache)
			totalChargeAmount, _, graduateStep := event_charge.ComputeMetricUsedChargePrice(metricUserUsedValue.UsedValue, metricMeteredCharge)
			meteredChargeStats = append(meteredChargeStats, &detail.UserMerchantMetricChargeStat{
				MetricId:          metricMeteredCharge.MetricId,
				Metric:            bean.SimplifyMerchantMetric(met),
				CurrentUsedValue:  metricUserUsedValue.UsedValue,
				MaxEventId:        metricUserUsedValue.MaxEventId,
				MinEventId:        metricUserUsedValue.MinEventId,
				ChargePricing:     metricMeteredCharge,
				TotalChargeAmount: totalChargeAmount,
				GraduatedStep:     graduateStep,
			})
		}
		for _, metricRecurringCharge := range planMetricBindingEntity.MetricRecurringCharge {
			met := query.GetMerchantMetric(ctx, metricRecurringCharge.MetricId)
			metricUserUsedValue := GetUserMetricCachedUseValue(ctx, merchantId, user.Id, met, one, reloadCache)
			totalChargeAmount, _, graduateStep := event_charge.ComputeMetricUsedChargePrice(metricUserUsedValue.UsedValue, metricRecurringCharge)
			recurringChargeStats = append(recurringChargeStats, &detail.UserMerchantMetricChargeStat{
				MetricId:          metricRecurringCharge.MetricId,
				Metric:            bean.SimplifyMerchantMetric(met),
				CurrentUsedValue:  metricUserUsedValue.UsedValue,
				MaxEventId:        metricUserUsedValue.MaxEventId,
				MinEventId:        metricUserUsedValue.MinEventId,
				ChargePricing:     metricRecurringCharge,
				TotalChargeAmount: totalChargeAmount,
				GraduatedStep:     graduateStep,
			})
		}
		return &detail.UserMetric{
			IsPaid:               one.Status == consts.SubStatusActive || one.Status == consts.SubStatusIncomplete,
			Product:              bean.SimplifyProduct(query.GetProductById(ctx, uint64(query.GetPlanById(ctx, one.PlanId).ProductId), merchantId)),
			User:                 bean.SimplifyUserAccount(user),
			Subscription:         bean.SimplifySubscription(ctx, one),
			Plan:                 bean.SimplifyPlan(query.GetPlanById(ctx, one.PlanId)),
			Addons:               addon2.GetSubscriptionAddonsByAddonJson(ctx, one.AddonData),
			LimitStats:           list,
			MeteredChargeStats:   meteredChargeStats,
			RecurringChargeStats: recurringChargeStats,
		}
	} else {
		return &detail.UserMetric{
			IsPaid:               false,
			User:                 bean.SimplifyUserAccount(user),
			Product:              nil,
			Subscription:         nil,
			Plan:                 nil,
			Addons:               nil,
			LimitStats:           list,
			MeteredChargeStats:   meteredChargeStats,
			RecurringChargeStats: recurringChargeStats,
		}
	}
}

func checkMetricUsedValue(ctx context.Context, merchantId uint64, user *entity.UserAccount, sub *entity.Subscription, met *entity.MerchantMetric, append int64) (int64, uint64, bool) {
	limitMap := GetUserMetricTotalLimits(ctx, merchantId, user.Id, sub)
	useValue := GetUserMetricCachedUseValue(ctx, merchantId, user.Id, met, sub, false)
	if metricLimit, ok := limitMap[met.Id]; ok {
		if met.AggregationType == metric.MetricAggregationTypeLatest || met.AggregationType == metric.MetricAggregationTypeMax {
			return useValue.UsedValue, metricLimit.TotalLimit, append <= int64(metricLimit.TotalLimit)
		} else {
			return useValue.UsedValue, metricLimit.TotalLimit, useValue.UsedValue+append <= int64(metricLimit.TotalLimit)
		}
	} else {
		// charge type
		return useValue.UsedValue, 0, true
	}
}

func GetUserMetricTotalLimits(ctx context.Context, merchantId uint64, userId uint64, sub *entity.Subscription) map[uint64]*detail.PlanMetricLimitDetail {
	var limitMap = make(map[uint64]*detail.PlanMetricLimitDetail)
	userSubPlans := user_sub_plan.UserSubPlanCachedListForMetric(ctx, merchantId, userId, sub, false)
	if len(userSubPlans) > 0 {
		g.Log().Infof(ctx, "GetUserMetricTotalLimits userId:%d subPlanId:%d userSubPlans:%s", userId, sub.PlanId, utility.MarshalToJsonString(userSubPlans))
		for _, subPlan := range userSubPlans {
			list := metric.MerchantMetricPlanLimitCachedList(ctx, merchantId, subPlan.PlanId, false)
			for _, planLimit := range list {
				if planLimit.Metric != nil {
					if _, ok := limitMap[planLimit.MetricId]; ok {
						limitMap[planLimit.MetricId].TotalLimit = limitMap[planLimit.MetricId].TotalLimit + planLimit.MetricLimit
						limitMap[planLimit.MetricId].PlanLimits = append(limitMap[planLimit.MetricId].PlanLimits, planLimit)
					} else {
						limitMap[planLimit.MetricId] = &detail.PlanMetricLimitDetail{
							MerchantId:          merchantId,
							UserId:              userId,
							MetricId:            planLimit.MetricId,
							Code:                planLimit.Metric.Code,
							MetricName:          planLimit.Metric.MetricName,
							Type:                planLimit.Metric.Type,
							AggregationType:     planLimit.Metric.AggregationType,
							AggregationProperty: planLimit.Metric.AggregationProperty,
							TotalLimit:          planLimit.MetricLimit,
							PlanLimits:          []*detail.MerchantMetricPlanLimitDetail{planLimit},
						}
					}
				}
			}
		}
	}
	return limitMap
}

const (
	UserMetricCacheKeyPrefix = "UserMetricCacheKeyV2Prefix_"
	UserMetricCacheKeyExpire = 15 * 24 * 60 * 60 // 15 days cache expire
)

type MetricUserUsedValue struct {
	UsedValue  int64  `json:"usedValue"`
	MaxEventId uint64 `json:"maxEventId"`
	MinEventId uint64 `json:"minEventId"`
}

func GetUserMetricCachedUseValue(ctx context.Context, merchantId uint64, userId uint64, met *entity.MerchantMetric, sub *entity.Subscription, reloadCache bool) (metricUserUsedValue *MetricUserUsedValue) {
	cacheKey := metricUserCacheKey(merchantId, userId, met, sub)
	if !reloadCache {
		get, err := g.Redis().Get(ctx, cacheKey)
		if err == nil && !get.IsNil() && !get.IsEmpty() {
			_ = utility.UnmarshalFromJsonString(get.String(), &metricUserUsedValue)
			if metricUserUsedValue != nil {
				return metricUserUsedValue
			}
		}
	}
	//var useValue uint64 = 0

	if merchantId > 0 {
		// count useValue from database
		if met.AggregationType == metric.MetricAggregationTypeLatest {
			var latestOne *entity.MerchantMetricEvent
			q := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				OrderDesc(dao.MerchantMetricEvent.Columns().GmtCreate)
			if met.Type == metric.MetricTypeLimitMetered {
				q = q.Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
					Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd)
			} else if met.Type == metric.MetricTypeChargeMetered {
				q = q.Where(dao.MerchantMetricEvent.Columns().ChargeStatus, 0)
			}
			err := q.Scan(&latestOne)
			utility.AssertError(err, "Server Error")
			if latestOne != nil {
				metricUserUsedValue = &MetricUserUsedValue{
					UsedValue:  latestOne.AggregationPropertyInt,
					MaxEventId: latestOne.Id,
					MinEventId: 0,
				}
			}
		} else if met.AggregationType == metric.MetricAggregationTypeMax {
			q := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0)
			if met.Type == metric.MetricTypeLimitMetered {
				q = q.Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
					Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd)
			} else if met.Type == metric.MetricTypeChargeMetered {
				q = q.Where(dao.MerchantMetricEvent.Columns().ChargeStatus, 0)
			}
			err := q.FieldMax(dao.MerchantMetricEvent.Columns().AggregationPropertyInt, "usedValue").
				FieldMax(dao.MerchantMetricEvent.Columns().Id, "maxEventId").
				FieldMin(dao.MerchantMetricEvent.Columns().Id, "minEventId").
				Scan(&metricUserUsedValue)
			utility.AssertError(err, "Server Error")
		} else {
			q := dao.MerchantMetricEvent.Ctx(ctx).
				Where(dao.MerchantMetricEvent.Columns().MerchantId, merchantId).
				Where(dao.MerchantMetricEvent.Columns().UserId, userId).
				Where(dao.MerchantMetricEvent.Columns().MetricId, int64(met.Id)).
				Where(dao.MerchantMetricEvent.Columns().SubscriptionIds, sub.SubscriptionId).
				Where(dao.MerchantMetricEvent.Columns().IsDeleted, 0)
			if met.Type == metric.MetricTypeLimitMetered {
				q = q.Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodStart, sub.CurrentPeriodStart).
					Where(dao.MerchantMetricEvent.Columns().SubscriptionPeriodEnd, sub.CurrentPeriodEnd)
			} else if met.Type == metric.MetricTypeChargeMetered {
				q = q.Where(dao.MerchantMetricEvent.Columns().ChargeStatus, 0)
			}
			err := q.FieldSum(dao.MerchantMetricEvent.Columns().AggregationPropertyInt, "usedValue").
				FieldMax(dao.MerchantMetricEvent.Columns().Id, "maxEventId").
				FieldMin(dao.MerchantMetricEvent.Columns().Id, "minEventId").
				Scan(&metricUserUsedValue)
			utility.AssertError(err, "Server Error")

		}
	}

	if metricUserUsedValue != nil {
		_, _ = g.Redis().Set(ctx, cacheKey, utility.MarshalToJsonString(metricUserUsedValue))
		_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)
	} else {
		metricUserUsedValue = &MetricUserUsedValue{}
	}

	return metricUserUsedValue
}

func appendMetricCachedUseValue(ctx context.Context, merchantId uint64, user *entity.UserAccount, met *entity.MerchantMetric, sub *entity.Subscription, append int64, maxEventId uint64) (metricUserUsedValue *MetricUserUsedValue) {
	cacheKey := metricUserCacheKey(merchantId, user.Id, met, sub)
	get, err := g.Redis().Get(ctx, cacheKey)
	if err == nil && !get.IsNil() && !get.IsEmpty() {
		_ = utility.UnmarshalFromJsonString(get.String(), &metricUserUsedValue)
	}
	if metricUserUsedValue != nil {
		newValue := metricUserUsedValue.UsedValue + append
		if met.AggregationType == metric.MetricAggregationTypeLatest {
			newValue = append
		} else if met.AggregationType == metric.MetricAggregationTypeMax {
			newValue = utility.MaxInt64(metricUserUsedValue.UsedValue, append)
		}
		metricUserUsedValue.UsedValue = newValue
		metricUserUsedValue.MaxEventId = maxEventId
		_, _ = g.Redis().Set(ctx, cacheKey, utility.MarshalToJsonString(metricUserUsedValue))
		_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)
		return metricUserUsedValue
	} else {
		return GetUserMetricCachedUseValue(ctx, merchantId, user.Id, met, sub, true)
	}
}

func metricUserCacheKey(merchantId uint64, userId uint64, met *entity.MerchantMetric, sub *entity.Subscription) string {
	cacheKey := fmt.Sprintf("%s_%d_%d_%d_%s_%d", UserMetricCacheKeyPrefix, merchantId, userId, met.Id, sub.SubscriptionId, sub.CurrentPeriodStart)
	return cacheKey
}

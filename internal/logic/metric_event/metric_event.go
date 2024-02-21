package metric_event

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee-api/internal/dao/oversea_pay"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/metric"
	"unibee-api/internal/logic/subscription/user_sub_plan"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

type MerchantMetricEventInternalReq struct {
	MerchantId       int64       `p:"merchantId" dc:"MerchantId" v:"required"`
	MetricCode       string      `p:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId   string      `p:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId  string      `p:"externalEventId" dc:"ExternalEventId, __unique__" v:"required"`
	MetricProperties *gjson.Json `p:"metricProperties" dc:"MetricProperties"`
}

func NewMerchantMetricEvent(ctx context.Context, req *MerchantMetricEventInternalReq) (*entity.MerchantMetricEvent, error) {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(len(req.MetricCode) > 0, "invalid metricCode")
	utility.Assert(len(req.ExternalEventId) > 0, "invalid externalEventId")
	utility.Assert(len(req.ExternalUserId) > 0, "invalid externalUserId")
	// user check
	user := query.GetUserAccountByExternalUserId(ctx, req.ExternalUserId)
	utility.Assert(user != nil, "user not found")
	// merchant check
	// metric check
	met := query.GetMerchantMetricByCode(ctx, req.MetricCode)
	utility.Assert(met != nil, "metric not found")
	utility.Assert(met.MerchantId == req.MerchantId, "code not match")

	// property determine
	var aggregationPropertyString = ""
	var aggregationPropertyInt uint64 = 1
	aggregationPropertyUniqueId := fmt.Sprintf("%d_%d_%d_%s", req.MerchantId, user.Id, met.Id, req.ExternalEventId)
	if met.AggregationType == metric.MetricAggregationTypeCount {
		// use aggregationPropertyInt, check properties
		aggregationPropertyInt = 1
	} else if met.AggregationType == metric.MetricAggregationTypeCountUnique {
		// use aggregationPropertyString, check properties
		utility.Assert(req.MetricProperties.Contains(met.AggregationProperty), fmt.Sprintf("property named '%s' not found in metricProperties json", met.AggregationProperty))
		// check value should be string
		utility.Assert(!req.MetricProperties.Get(met.AggregationProperty).IsMap(), fmt.Sprintf("property named '%s' is not string type, it's Map", met.AggregationProperty))
		utility.Assert(!req.MetricProperties.Get(met.AggregationProperty).IsFloat(), fmt.Sprintf("property named '%s' is not string type, it's Float", met.AggregationProperty))
		utility.Assert(!req.MetricProperties.Get(met.AggregationProperty).IsStruct(), fmt.Sprintf("property named '%s' is not string type, it's Struct", met.AggregationProperty))
		utility.Assert(!req.MetricProperties.Get(met.AggregationProperty).IsSlice(), fmt.Sprintf("property named '%s' is not string type, it's Slice", met.AggregationProperty))
		utility.Assert(!req.MetricProperties.Get(met.AggregationProperty).IsEmpty(), fmt.Sprintf("property named '%s' is not string type, it's Empty", met.AggregationProperty))
		utility.Assert(!req.MetricProperties.Get(met.AggregationProperty).IsNil(), fmt.Sprintf("property named '%s' is not string type, it's Empty", met.AggregationProperty))
		utility.Assert(!req.MetricProperties.Get(met.AggregationProperty).IsUint(), fmt.Sprintf("property named '%s' is not string type, it's Uint", met.AggregationProperty))
		utility.Assert(!req.MetricProperties.Get(met.AggregationProperty).IsInt(), fmt.Sprintf("property named '%s' is not string type, it's Int", met.AggregationProperty))
		aggregationPropertyString = req.MetricProperties.Get(met.AggregationProperty).String()
		aggregationPropertyInt = 1
		// count unique should replace uniqueId eventId with unique property
		aggregationPropertyUniqueId = fmt.Sprintf("%d_%d_%d_%s", req.MerchantId, met.Id, user.Id, aggregationPropertyString)
	} else {
		// use aggregationPropertyInt, check properties
		utility.Assert(req.MetricProperties.Contains(met.AggregationProperty), fmt.Sprintf("property named '%s' not found in metricProperties json", met.AggregationProperty))
		// check value should be int
		utility.Assert(req.MetricProperties.Get(met.AggregationProperty).IsUint(), fmt.Sprintf("property named '%s' is not Uint type", met.AggregationProperty))
		aggregationPropertyInt = req.MetricProperties.Get(met.AggregationProperty).Uint64()
	}

	if met.Type == metric.MetricTypeLimitMetered {
		// need check if metric limit reached and reject it
		useValue, metricLimit, check := checkMetricLimitReached(ctx, req.MerchantId, user, met, aggregationPropertyInt)
		utility.Assert(check, fmt.Sprintf("metric limit reached, current use: %d, limit: %d", useValue, metricLimit))
	}

	var one *entity.MerchantMetricEvent
	err := dao.MerchantMetricPlanLimit.Ctx(ctx).
		Where(entity.MerchantMetricEvent{AggregationPropertyUniqueId: aggregationPropertyUniqueId}).
		Scan(&one)
	utility.AssertError(err, "server error")
	utility.Assert(one == nil, "same event with externalEventId or uniqueProperty exist")

	one = &entity.MerchantMetricEvent{
		MerchantId:                  req.MerchantId,
		MetricId:                    int64(met.Id),
		ExternalEventId:             req.ExternalEventId,
		UserId:                      int64(user.Id),
		AggregationPropertyData:     req.MetricProperties.String(),
		AggregationPropertyInt:      int64(aggregationPropertyInt),
		AggregationPropertyString:   aggregationPropertyString,
		AggregationPropertyUniqueId: aggregationPropertyUniqueId,
		SubscriptionIds:             "", //todo mark
		SubscriptionPeriodEnd:       0,  //todo mark
		SubscriptionPeriodStart:     0,  //todo mark
		CreateTime:                  gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantMetricEvent.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Errorf(ctx, "event insert err:%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "event server error", nil))
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	if met.Type == metric.MetricTypeLimitMetered {
		// append the metric limit usage

	}

	return one, nil
}

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
		return useValue, metricLimit.TotalLimit, useValue+append <= metricLimit.TotalLimit
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
			useValue = 0
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

func appendMetricLimitCachedUseValue(ctx context.Context, user *entity.UserAccount, merchantId int64, met *entity.MerchantMetric, append uint64) {
	cacheKey := fmt.Sprintf("%s_%d_%d_%d", UserMetricCacheKeyPrefix, merchantId, user.Id, met.Id)
	get, err := g.Redis().Get(ctx, cacheKey)
	if err == nil && !get.IsNil() && !get.IsEmpty() && (get.IsUint() || get.IsInt()) {
		_, _ = g.Redis().Set(ctx, cacheKey, get.Uint64()+append)
		_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)
	} else {
		_, _ = g.Redis().Set(ctx, cacheKey, append)
		_, _ = g.Redis().Expire(ctx, cacheKey, UserMetricCacheKeyExpire)
	}
}

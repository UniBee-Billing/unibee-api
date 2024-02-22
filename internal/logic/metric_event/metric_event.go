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
	"unibee-api/internal/logic/metric"
	entity "unibee-api/internal/model/entity/oversea_pay"
	"unibee-api/internal/query"
	"unibee-api/utility"
)

type MerchantMetricEventInternalReq struct {
	MerchantId       uint64      `p:"merchantId" dc:"MerchantId" v:"required"`
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
	user := query.GetUserAccountByExternalUserId(ctx, req.MerchantId, req.ExternalUserId)
	utility.Assert(user != nil, "user not found")
	// merchant check
	// metric check
	met := query.GetMerchantMetricByCode(ctx, req.MetricCode)
	utility.Assert(met != nil, "metric not found")
	utility.Assert(met.MerchantId == req.MerchantId, "code not match")
	// check the only subscription, todo mark limit add subscription and cycle reset metric support
	sub := query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, int64(user.Id), req.MerchantId)
	utility.Assert(sub != nil, "user has no subscription")

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
		Where(dao.MerchantMetricEvent.Columns().AggregationPropertyUniqueId, aggregationPropertyUniqueId).
		Scan(&one)
	utility.AssertError(err, "server error")
	utility.Assert(one == nil, "same event with externalEventId or uniqueProperty exist")

	one = &entity.MerchantMetricEvent{
		MerchantId:                  req.MerchantId,
		MetricId:                    int64(met.Id),
		ExternalEventId:             req.ExternalEventId,
		UserId:                      int64(user.Id),
		AggregationPropertyData:     req.MetricProperties.String(),
		AggregationPropertyInt:      aggregationPropertyInt,
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
		appendMetricLimitCachedUseValue(ctx, req.MerchantId, user, met, aggregationPropertyInt)
	}

	return one, nil
}

func DelMerchantMetricEvent(ctx context.Context, req *MerchantMetricEventInternalReq) error {
	utility.Assert(req.MerchantId > 0, "invalid merchantId")
	utility.Assert(len(req.MetricCode) > 0, "invalid metricCode")
	utility.Assert(len(req.ExternalEventId) > 0, "invalid externalEventId")
	utility.Assert(len(req.ExternalUserId) > 0, "invalid externalUserId")
	// user check
	user := query.GetUserAccountByExternalUserId(ctx, req.MerchantId, req.ExternalUserId)
	utility.Assert(user != nil, "user not found")
	// merchant check
	// metric check
	met := query.GetMerchantMetricByCode(ctx, req.MetricCode)
	utility.Assert(met != nil, "metric not found")
	utility.Assert(met.MerchantId == req.MerchantId, "code not match")
	var list []*entity.MerchantMetricEvent
	err := dao.MerchantMetricEvent.Ctx(ctx).
		Where(dao.MerchantMetricEvent.Columns().MerchantId, req.MerchantId).
		Where(dao.MerchantMetricEvent.Columns().MetricId, met.MerchantId).
		Where(dao.MerchantMetricEvent.Columns().UserId, int64(user.Id)).
		Where(dao.MerchantMetricEvent.Columns().ExternalEventId, req.ExternalEventId).
		Scan(&list)
	if err != nil {
		return err
	}
	utility.Assert(len(list) == 1, "event not found")
	_, err = dao.MerchantMetricEvent.Ctx(ctx).Data(g.Map{
		dao.MerchantMetricEvent.Columns().IsDeleted: gtime.Now().Timestamp(),
		dao.MerchantMetricEvent.Columns().GmtModify: gtime.Now(),
	}).Where(dao.MerchantMetricEvent.Columns().Id, list[0].Id).OmitNil().Update()
	if err != nil {
		return err
	}
	return nil
}

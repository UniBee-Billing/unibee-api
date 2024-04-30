package metric_event

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee/internal/dao/oversea_pay"
	"unibee/internal/logic/metric"
	entity "unibee/internal/model/entity/oversea_pay"
	"unibee/internal/query"
	"unibee/utility"
)

type MerchantMetricEventInternalReq struct {
	MerchantId       uint64      `json:"merchantId" dc:"MerchantId" v:"required"`
	MetricCode       string      `json:"metricCode" dc:"MetricCode" v:"required"`
	ExternalUserId   string      `json:"externalUserId" dc:"ExternalUserId" v:"required"`
	ExternalEventId  string      `json:"externalEventId" dc:"ExternalEventId, __unique__" v:"required"`
	MetricProperties *gjson.Json `json:"metricProperties" dc:"MetricProperties"`
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
	// check the only paid subscription
	sub := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, user.Id, req.MerchantId)
	utility.Assert(sub != nil, "user has no active subscription")

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

	var metricLimit uint64 = 0
	var usedValue uint64 = 0
	var check = false
	if met.Type == metric.MetricTypeLimitMetered {
		// need check if metric limit reached and reject it
		usedValue, metricLimit, check = checkMetricLimitReached(ctx, req.MerchantId, user, sub, met, aggregationPropertyInt)
		utility.Assert(check, fmt.Sprintf("metric limit reached, current used: %d, limit: %d", usedValue, metricLimit))
	}

	var one *entity.MerchantMetricEvent
	err := dao.MerchantMetricEvent.Ctx(ctx).
		Where(dao.MerchantMetricEvent.Columns().AggregationPropertyUniqueId, aggregationPropertyUniqueId).
		Scan(&one)
	utility.AssertError(err, "Server Error")
	//if one != nil && one.AggregationPropertyInt == aggregationPropertyInt {
	//	// same event
	//	return one, nil
	//} else {
	utility.Assert(one == nil, "same event with externalEventId or uniqueProperty exist")
	//}

	one = &entity.MerchantMetricEvent{
		MerchantId:                  req.MerchantId,
		MetricId:                    met.Id,
		ExternalEventId:             req.ExternalEventId,
		UserId:                      int64(user.Id),
		AggregationPropertyData:     req.MetricProperties.String(),
		AggregationPropertyInt:      aggregationPropertyInt,
		AggregationPropertyString:   aggregationPropertyString,
		AggregationPropertyUniqueId: aggregationPropertyUniqueId,
		SubscriptionIds:             sub.SubscriptionId,
		SubscriptionPeriodStart:     sub.CurrentPeriodStart,
		SubscriptionPeriodEnd:       sub.CurrentPeriodEnd,
		CreateTime:                  gtime.Now().Timestamp(),
		MetricLimit:                 metricLimit,
		Used:                        usedValue,
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
		usedValue = appendMetricLimitCachedUseValue(ctx, req.MerchantId, user, met, sub, aggregationPropertyInt)
	}
	one.Used = usedValue

	go func() {
		// update background
		backgroundCtx := context.Background()
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				if err != nil {
					g.Log().Errorf(backgroundCtx, "NewMerchantMetricEvent Update UsedValue panic error:%s", err.Error())
				} else {
					g.Log().Errorf(backgroundCtx, "NewMerchantMetricEvent Update UsedValue panic error:%s", err)
				}
				return
			}
		}()
		_, err = dao.MerchantMetricEvent.Ctx(backgroundCtx).Data(g.Map{
			dao.MerchantMetricEvent.Columns().Used:      usedValue,
			dao.MerchantMetricEvent.Columns().GmtModify: gtime.Now(),
		}).Where(dao.MerchantMetricEvent.Columns().Id, one.Id).OmitNil().Update()
		if err != nil {
			g.Log().Errorf(backgroundCtx, "NewMerchantMetricEvent Update UsedValue err:%s", err.Error())
		}
	}()

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

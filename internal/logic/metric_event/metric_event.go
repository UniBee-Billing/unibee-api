package metric_event

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	dao "unibee-api/internal/dao/oversea_pay"
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
	utility.Assert(len(req.MetricCode) > 0, "MetricCode is nil")
	// user check
	user := query.GetUserAccountByExternalUserId(ctx, req.ExternalUserId)
	utility.Assert(user != nil, "user not found")
	// merchant check
	// metric check
	oneMetric := query.GetMerchantMetricByCode(ctx, req.MetricCode)
	utility.Assert(oneMetric != nil, "metric not found")
	var one *entity.MerchantMetricEvent
	err := dao.MerchantMetricPlanLimit.Ctx(ctx).
		Where(entity.MerchantMetricEvent{MerchantId: req.MerchantId}).
		Where(entity.MerchantMetricEvent{MetricId: int64(oneMetric.Id)}).
		Where(entity.MerchantMetricEvent{ExternalEventId: req.ExternalEventId}).
		Where(entity.MerchantMetricEvent{IsDeleted: 0}).
		Scan(&one)
	utility.AssertError(err, "server error")
	utility.Assert(one == nil, "externalEventId exist")
	one = &entity.MerchantMetricEvent{
		MerchantId:                  req.MerchantId,
		MetricId:                    int64(oneMetric.Id),
		ExternalEventId:             req.ExternalEventId,
		UserId:                      int64(user.Id),
		AggregationPropertyData:     req.MetricProperties.String(),
		AggregationPropertyInt:      0,  //todo mark
		AggregationPropertyString:   "", //todo mark
		AggregationPropertyUniqueId: "", // todo mark
		SubscriptionIds:             "", //todo mark
		SubscriptionPeriodEnd:       0,  //todo mark
		SubscriptionPeriodStart:     0,  //todo mark
		CreateTime:                  gtime.Now().Timestamp(),
	}
	result, err := dao.MerchantMetricEvent.Ctx(ctx).Data(one).OmitNil().Insert(one)
	if err != nil {
		g.Log().Errorf(ctx, "NewMerchantMetricEvent Insert err:%s", err.Error())
		return nil, gerror.NewCode(gcode.New(500, "server error", nil))
	}
	id, _ := result.LastInsertId()
	one.Id = uint64(id)

	return one, nil
}

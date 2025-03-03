package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/bean"
	"unibee/api/merchant/metric"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/metric_event"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func (c *ControllerMetric) NewEvent(ctx context.Context, req *metric.NewEventReq) (res *metric.NewEventRes, err error) {
	var one *entity.UserAccount
	if req.UserId > 0 {
		one = query.GetUserAccountById(ctx, req.UserId)
	} else if len(req.Email) > 0 {
		one = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	} else if len(req.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
	}
	if one == nil {
		return nil, gerror.New("user not found, should provides one of three options, UserId, ExternalUserId, or Email")
	}
	event, err := metric_event.NewMerchantMetricEvent(ctx, &metric_event.MerchantMetricEventInternalReq{
		MerchantId:          _interface.GetMerchantId(ctx),
		MetricCode:          req.MetricCode,
		UserId:              one.Id,
		ExternalEventId:     req.ExternalEventId,
		MetricProperties:    req.MetricProperties,
		ProductId:           req.ProductId,
		AggregationValue:    req.AggregationValue,
		AggregationUniqueId: req.AggregationUniqueId,
	})
	if err != nil {
		return nil, err
	}
	return &metric.NewEventRes{MerchantMetricEvent: bean.SimplifyMerchantMetricEvent(event)}, nil
}

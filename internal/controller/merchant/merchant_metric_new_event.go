package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/merchant/metric"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/metric_event"
)

func (c *ControllerMetric) NewEvent(ctx context.Context, req *metric.NewEventReq) (res *metric.NewEventRes, err error) {
	event, err := metric_event.NewMerchantMetricEvent(ctx, &metric_event.MerchantMetricEventInternalReq{
		MerchantId:       _interface.GetMerchantId(ctx),
		MetricCode:       req.MetricCode,
		ExternalUserId:   req.ExternalUserId,
		ExternalEventId:  req.ExternalEventId,
		MetricProperties: req.MetricProperties,
	})
	if err != nil {
		return nil, err
	}
	return &metric.NewEventRes{MerchantMetricEvent: bean.SimplifyMerchantMetricEvent(event)}, nil
}

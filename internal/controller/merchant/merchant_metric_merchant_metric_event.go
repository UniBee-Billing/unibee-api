package merchant

import (
	"context"
	"unibee-api/api/merchant/metric"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/metric_event"
)

func (c *ControllerMetric) MerchantMetricEvent(ctx context.Context, req *metric.MerchantMetricEventReq) (res *metric.MerchantMetricEventRes, err error) {
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
	return &metric.MerchantMetricEventRes{MerchantMetricEvent: event}, nil
}
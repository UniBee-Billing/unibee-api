package merchant

import (
	"context"
	"unibee-api/api/merchant/metric"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/metric_event"
)

func (c *ControllerMetric) DelMerchantMetricEvent(ctx context.Context, req *metric.DelMerchantMetricEventReq) (res *metric.DelMerchantMetricEventRes, err error) {
	err = metric_event.DelMerchantMetricEvent(ctx, &metric_event.MerchantMetricEventInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		MetricCode:      req.MetricCode,
		ExternalUserId:  req.ExternalUserId,
		ExternalEventId: req.ExternalEventId,
	})
	if err != nil {
		return nil, err
	}
	return &metric.DelMerchantMetricEventRes{}, nil
}

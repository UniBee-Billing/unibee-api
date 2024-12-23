package merchant

import (
	"context"
	"unibee/api/merchant/metric"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/metric_event"
)

func (c *ControllerMetric) DeleteEvent(ctx context.Context, req *metric.DeleteEventReq) (res *metric.DeleteEventRes, err error) {
	err = metric_event.DelMerchantMetricEvent(ctx, &metric_event.MerchantMetricEventInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		MetricCode:      req.MetricCode,
		ExternalUserId:  req.ExternalUserId,
		ExternalEventId: req.ExternalEventId,
	})
	if err != nil {
		return nil, err
	}
	return &metric.DeleteEventRes{}, nil
}

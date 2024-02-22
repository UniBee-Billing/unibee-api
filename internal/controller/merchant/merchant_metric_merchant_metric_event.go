package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/metric_event"
	"unibee-api/utility"

	"unibee-api/api/merchant/metric"
)

func (c *ControllerMetric) MerchantMetricEvent(ctx context.Context, req *metric.MerchantMetricEventReq) (res *metric.MerchantMetricEventRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
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

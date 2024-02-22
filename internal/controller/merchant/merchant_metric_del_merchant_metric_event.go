package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/metric_event"
	"unibee-api/utility"

	"unibee-api/api/merchant/metric"
)

func (c *ControllerMetric) DelMerchantMetricEvent(ctx context.Context, req *metric.DelMerchantMetricEventReq) (res *metric.DelMerchantMetricEventRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//Merchant User Check
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
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

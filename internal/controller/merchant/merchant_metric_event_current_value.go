package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/metric_event"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) EventCurrentValue(ctx context.Context, req *metric.EventCurrentValueReq) (res *metric.EventCurrentValueRes, err error) {
	var one *entity.UserAccount
	if req.UserId > 0 {
		one = query.GetUserAccountById(ctx, req.UserId)
	} else if len(req.Email) > 0 {
		one = query.GetUserAccountByEmail(ctx, _interface.GetMerchantId(ctx), req.Email)
	} else if len(req.ExternalUserId) > 0 {
		one = query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
	}
	utility.Assert(one != nil, "user not found, should provides one of three options, UserId, ExternalUserId, or Email")
	value := metric_event.MerchantMetricEventCurrentValue(ctx, &metric_event.MerchantMetricEventInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		MetricCode: req.MetricCode,
		UserId:     one.Id,
		ProductId:  req.ProductId,
	})
	return &metric.EventCurrentValueRes{CurrentValue: value}, nil
}

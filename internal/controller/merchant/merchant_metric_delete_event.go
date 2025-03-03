package merchant

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"unibee/api/merchant/metric"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/metric_event"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
)

func (c *ControllerMetric) DeleteEvent(ctx context.Context, req *metric.DeleteEventReq) (res *metric.DeleteEventRes, err error) {
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
	err = metric_event.DelMerchantMetricEvent(ctx, &metric_event.MerchantMetricEventInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		MetricCode:      req.MetricCode,
		UserId:          one.Id,
		ExternalEventId: req.ExternalEventId,
	})
	if err != nil {
		return nil, err
	}
	return &metric.DeleteEventRes{}, nil
}

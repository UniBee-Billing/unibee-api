package user

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/metric_event"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/metric"
)

func (c *ControllerMetric) UserMetric(ctx context.Context, req *metric.UserMetricReq) (res *metric.UserMetricRes, err error) {
	user := query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
	utility.Assert(user != nil, "user not found")
	return &metric.UserMetricRes{UserMetric: metric_event.GetUserMetricStat(ctx, _interface.GetMerchantId(ctx), user, req.ProductId, false)}, nil
}

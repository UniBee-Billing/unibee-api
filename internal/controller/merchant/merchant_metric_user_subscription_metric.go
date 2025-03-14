package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/metric_event"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/metric"
)

func (c *ControllerMetric) UserSubscriptionMetric(ctx context.Context, req *metric.UserSubscriptionMetricReq) (res *metric.UserSubscriptionMetricRes, err error) {
	utility.Assert(len(req.SubscriptionId) > 0, "subscription id should not be empty")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	user := query.GetUserAccountById(ctx, sub.UserId)
	utility.Assert(user != nil, "user not found")
	return &metric.UserSubscriptionMetricRes{UserMetric: metric_event.GetUserSubscriptionMetricStat(ctx, _interface.GetMerchantId(ctx), user, sub, false)}, nil
}

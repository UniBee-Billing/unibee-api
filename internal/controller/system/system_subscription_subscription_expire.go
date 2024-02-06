package system

import (
	"context"
	"unibee-api/internal/cronjob/sub"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"unibee-api/api/system/subscription"
)

func (c *ControllerSubscription) SubscriptionExpire(ctx context.Context, req *subscription.SubscriptionExpireReq) (res *subscription.SubscriptionExpireRes, err error) {
	one := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(one != nil, "sub not found")
	err = sub.SubscriptionExpire(ctx, one, "AdminExpire")
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionExpireRes{}, nil
}

package system

import (
	"context"
	"unibee/internal/logic/subscription/billingcycle/expire"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/system/subscription"
)

func (c *ControllerSubscription) SubscriptionExpire(ctx context.Context, req *subscription.SubscriptionExpireReq) (res *subscription.SubscriptionExpireRes, err error) {
	one := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(one != nil, "sub not found")
	err = expire.SubscriptionExpire(ctx, one, "AdminExpire")
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionExpireRes{}, nil
}

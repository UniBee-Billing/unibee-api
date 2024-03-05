package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) CancelLastCancelAtPeriodEnd(ctx context.Context, req *subscription.CancelLastCancelAtPeriodEndReq) (res *subscription.CancelLastCancelAtPeriodEndRes, err error) {
	err = service.SubscriptionCancelLastCancelAtPeriodEnd(ctx, req.SubscriptionId, false)
	if err != nil {
		return nil, err
	}
	return &subscription.CancelLastCancelAtPeriodEndRes{}, nil
}

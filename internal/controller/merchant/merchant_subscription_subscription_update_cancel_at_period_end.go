package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) CancelAtPeriodEnd(ctx context.Context, req *subscription.CancelAtPeriodEndReq) (res *subscription.CancelAtPeriodEndRes, err error) {
	err = service.SubscriptionCancelAtPeriodEnd(ctx, req.SubscriptionId, false, -1)
	if err != nil {
		return nil, err
	}
	return &subscription.CancelAtPeriodEndRes{}, nil
}

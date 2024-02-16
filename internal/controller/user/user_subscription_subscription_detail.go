package user

import (
	"context"
	"unibee-api/api/user/subscription"
	"unibee-api/internal/logic/subscription/service"
)

func (c *ControllerSubscription) SubscriptionDetail(ctx context.Context, req *subscription.SubscriptionDetailReq) (res *subscription.SubscriptionDetailRes, err error) {
	return service.SubscriptionDetail(ctx, req.SubscriptionId)
}

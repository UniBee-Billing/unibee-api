package user

import (
	"context"
	"unibee/api/user/subscription"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) SubscriptionCreate(ctx context.Context, req *subscription.SubscriptionCreateReq) (res *subscription.SubscriptionCreateRes, err error) {
	
	return service.SubscriptionCreate(ctx, req)
}

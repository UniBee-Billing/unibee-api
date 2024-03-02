package user

import (
	"context"
	"unibee/api/user/subscription"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error) {

	return service.SubscriptionCreate(ctx, req)
}

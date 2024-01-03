package user

import (
	"context"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerSubscription) SubscriptionCreate(ctx context.Context, req *subscription.SubscriptionCreateReq) (res *subscription.SubscriptionCreateRes, err error) {
	one, err := service.SubscriptionCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionCreateRes{Subscription: one}, nil
}

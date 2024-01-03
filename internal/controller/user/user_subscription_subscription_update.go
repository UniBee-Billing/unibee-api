package user

import (
	"context"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerSubscription) SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq) (res *subscription.SubscriptionUpdateRes, err error) {
	one, err := service.SubscriptionUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdateRes{SubscriptionPendingUpdate: one}, nil
}

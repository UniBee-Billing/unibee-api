package user

import (
	"context"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/logic/subscription/service"
)

func (c *ControllerSubscription) SubscriptionCancel(ctx context.Context, req *subscription.SubscriptionCancelReq) (res *subscription.SubscriptionCancelRes, err error) {
	err = service.SubscriptionCancel(ctx, req.SubscriptionId)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionCancelRes{}, nil
}

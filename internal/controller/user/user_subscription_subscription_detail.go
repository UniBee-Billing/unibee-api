package user

import (
	"context"
	"unibee-api/api/user/subscription"
	"unibee-api/internal/logic/subscription/service"
)

func (c *ControllerSubscription) SubscriptionDetail(ctx context.Context, req *subscription.SubscriptionDetailReq) (res *subscription.SubscriptionDetailRes, err error) {
	detail, err := service.SubscriptionDetail(ctx, req.SubscriptionId)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionDetailRes{
		User:                                detail.User,
		Subscription:                        detail.Subscription,
		Plan:                                detail.Plan,
		Gateway:                             detail.Gateway,
		Addons:                              detail.Addons,
		UnfinishedSubscriptionPendingUpdate: detail.UnfinishedSubscriptionPendingUpdate,
	}, nil
}

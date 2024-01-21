package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/merchant/subscription"
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
		Channel:                             detail.Channel,
		Addons:                              detail.Addons,
		UnfinishedSubscriptionPendingUpdate: detail.UnfinishedSubscriptionPendingUpdate,
	}, nil
}

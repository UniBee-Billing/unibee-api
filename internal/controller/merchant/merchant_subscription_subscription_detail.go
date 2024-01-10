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
		Subscription:               detail.Subscription,
		Plan:                       detail.Plan,
		Channel:                    detail.Channel,
		Addons:                     detail.Addons,
		SubscriptionPendingUpdates: detail.SubscriptionPendingUpdates,
	}, nil
}

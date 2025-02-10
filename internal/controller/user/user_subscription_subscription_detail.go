package user

import (
	"context"
	"unibee/api/user/subscription"
	"unibee/internal/logic/subscription/service/detail"
)

func (c *ControllerSubscription) Detail(ctx context.Context, req *subscription.DetailReq) (res *subscription.DetailRes, err error) {
	detail, err := detail.SubscriptionDetail(ctx, req.SubscriptionId)
	if err != nil {
		return nil, err
	}
	return &subscription.DetailRes{
		User:                                detail.User,
		Subscription:                        detail.Subscription,
		Plan:                                detail.Plan,
		Gateway:                             detail.Gateway,
		AddonParams:                         detail.AddonParams,
		Addons:                              detail.Addons,
		LatestInvoice:                       detail.LatestInvoice,
		UnfinishedSubscriptionPendingUpdate: detail.UnfinishedSubscriptionPendingUpdate,
	}, nil
}

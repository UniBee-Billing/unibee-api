package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/service/detail"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Detail(ctx context.Context, req *subscription.DetailReq) (res *subscription.DetailRes, err error) {
	detail, err := detail.SubscriptionDetail(ctx, req.SubscriptionId)
	if err != nil {
		return nil, err
	}
	utility.Assert(detail.Subscription.MerchantId == _interface.GetMerchantId(ctx), "wrong merchant account")
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

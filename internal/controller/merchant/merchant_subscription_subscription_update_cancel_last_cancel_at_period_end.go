package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) CancelLastCancelAtPeriodEnd(ctx context.Context, req *subscription.CancelLastCancelAtPeriodEndReq) (res *subscription.CancelLastCancelAtPeriodEndRes, err error) {
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		one := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx), req.ProductId)
		if one != nil {
			req.SubscriptionId = one.SubscriptionId
		} else {
			return &subscription.CancelLastCancelAtPeriodEndRes{}, nil
		}
	}

	err = service.SubscriptionCancelLastCancelAtPeriodEnd(ctx, req.SubscriptionId, false)
	if err != nil {
		return nil, err
	}
	return &subscription.CancelLastCancelAtPeriodEndRes{}, nil
}

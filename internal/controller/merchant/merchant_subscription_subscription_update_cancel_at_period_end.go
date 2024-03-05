package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) CancelAtPeriodEnd(ctx context.Context, req *subscription.CancelAtPeriodEndReq) (res *subscription.CancelAtPeriodEndRes, err error) {
	err = service.SubscriptionCancelAtPeriodEnd(ctx, req.SubscriptionId, false, int64(_interface.BizCtx().Get(ctx).MerchantMember.Id))
	if err != nil {
		return nil, err
	}
	return &subscription.CancelAtPeriodEndRes{}, nil
}

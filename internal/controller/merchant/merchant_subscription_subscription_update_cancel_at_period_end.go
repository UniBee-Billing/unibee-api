package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/subscription/service"
	"unibee-api/utility"

	"unibee-api/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdateCancelAtPeriodEnd(ctx context.Context, req *subscription.SubscriptionUpdateCancelAtPeriodEndReq) (res *subscription.SubscriptionUpdateCancelAtPeriodEndRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}

	err = service.SubscriptionCancelAtPeriodEnd(ctx, req.SubscriptionId, false, int64(_interface.BizCtx().Get(ctx).MerchantUser.Id))
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdateCancelAtPeriodEndRes{}, nil
}

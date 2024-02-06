package user

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/subscription/service"
	"unibee-api/utility"

	"unibee-api/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdateCancelLastCancelAtPeriodEnd(ctx context.Context, req *subscription.SubscriptionUpdateCancelLastCancelAtPeriodEndReq) (res *subscription.SubscriptionUpdateCancelLastCancelAtPeriodEndRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		//utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == sub.UserId, "userId not match") // todo mark
	}
	err = service.SubscriptionCancelLastCancelAtPeriodEnd(ctx, req.SubscriptionId, false)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdateCancelLastCancelAtPeriodEndRes{}, nil
}

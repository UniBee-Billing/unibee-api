package user

import (
	"context"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) CancelAtPeriodEnd(ctx context.Context, req *subscription.CancelAtPeriodEndReq) (res *subscription.CancelAtPeriodEndRes, err error) {
	if !config.GetConfigInstance().IsLocal() {
		utility.Assert(_interface.Context().Get(ctx).User != nil, "auth failure,not login")
		//utility.Assert(int64(_interface.Context().Get(ctx).User.Id) == sub.UserId, "userId not match") // todo mark
	}
	err = service.SubscriptionCancelAtPeriodEnd(ctx, req.SubscriptionId, false, 0)
	if err != nil {
		return nil, err
	}
	return &subscription.CancelAtPeriodEndRes{}, nil
}

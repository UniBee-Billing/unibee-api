package user

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionCancel(ctx context.Context, req *subscription.SubscriptionCancelReq) (res *subscription.SubscriptionCancelRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).User.Id > 0, "userId invalid")
	}

	utility.Assert(len(req.SubscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.UserId == int64(_interface.BizCtx().Get(ctx).User.Id), "no permission")
	utility.Assert(sub.Status != consts.SubStatusCancelled, "subscription already cancelled")
	utility.Assert(sub.Status == consts.SubStatusCreate, "subscription not in create status")

	err = service.SubscriptionCancel(ctx, req.SubscriptionId, false, false, "User Cancel Create Subscription")
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionCancelRes{}, nil
}

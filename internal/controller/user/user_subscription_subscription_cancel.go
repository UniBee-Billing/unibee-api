package user

import (
	"context"
	"unibee/internal/cmd/config"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) Cancel(ctx context.Context, req *subscription.CancelReq) (res *subscription.CancelRes, err error) {
	if !config.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.Context().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).User.Id > 0, "userId invalid")
	}

	utility.Assert(len(req.SubscriptionId) > 0, "subscriptionId not found")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)
	utility.Assert(sub != nil, "subscription not found")
	utility.Assert(sub.UserId == _interface.Context().Get(ctx).User.Id, "no permission")
	utility.Assert(sub.Status != consts.SubStatusCancelled, "subscription already cancelled")
	utility.Assert(sub.Status == consts.SubStatusPending || sub.Status == consts.SubStatusProcessing, "subscription not in pending or processing status")

	err = service.SubscriptionCancel(ctx, req.SubscriptionId, false, false, "CancelledByUser")
	if err != nil {
		return nil, err
	}
	return &subscription.CancelRes{}, nil
}

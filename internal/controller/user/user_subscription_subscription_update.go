package user

import (
	"context"
	"unibee/api/user/subscription"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	utility.Assert(req.WithImmediateEffect == 0, "WithImmediateEffect not support in user_portal")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)

	//Update 可以由 Admin 操作，service 层不做用户校验
	if !config.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == sub.UserId, "userId not match")
	}
	return service.SubscriptionUpdate(ctx, req, 0)
}

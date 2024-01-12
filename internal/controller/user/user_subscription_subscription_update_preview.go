package user

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq) (res *subscription.SubscriptionUpdatePreviewRes, err error) {
	utility.Assert(req != nil, "req not found")
	utility.Assert(req.NewPlanId > 0, "PlanId invalid")
	utility.Assert(len(req.SubscriptionId) > 0, "SubscriptionId invalid")
	utility.Assert(req.WithImmediateEffect == 0, "WithImmediateEffect not support in user_portal")
	sub := query.GetSubscriptionBySubscriptionId(ctx, req.SubscriptionId)

	//Update 可以由 Admin 操作，service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == sub.UserId, "userId not match")
	}
	prepare, err := service.SubscriptionUpdatePreview(ctx, req, 0, 0)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdatePreviewRes{
		TotalAmount:   prepare.TotalAmount,
		Currency:      prepare.Currency,
		Invoice:       prepare.Invoice,
		ProrationDate: prepare.ProrationDate,
	}, nil
}

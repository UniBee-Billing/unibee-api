package merchant

import (
	"context"
	subscription2 "unibee/api/user/subscription"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq) (res *subscription.SubscriptionUpdateRes, err error) {

	var merchantMemberId int64
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
		merchantMemberId = int64(_interface.BizCtx().Get(ctx).MerchantMember.Id)
	}
	update, err := service.SubscriptionUpdate(ctx, &subscription2.SubscriptionUpdateReq{
		SubscriptionId:      req.SubscriptionId,
		NewPlanId:           req.NewPlanId,
		Quantity:            req.Quantity,
		AddonParams:         req.AddonParams,
		WithImmediateEffect: req.WithImmediateEffect,
		ConfirmTotalAmount:  req.ConfirmTotalAmount,
		ConfirmCurrency:     req.ConfirmCurrency,
		ProrationDate:       req.ProrationDate,
	}, merchantMemberId)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdateRes{
		SubscriptionPendingUpdate: update.SubscriptionPendingUpdate,
		Paid:                      update.Paid,
		Link:                      update.Link,
	}, nil
}

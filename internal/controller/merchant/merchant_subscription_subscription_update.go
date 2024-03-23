package merchant

import (
	"context"
	subscription2 "unibee/api/user/subscription"
	"unibee/internal/cmd/config"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error) {

	var merchantMemberId int64
	if !config.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.Context().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.Context().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
		merchantMemberId = int64(_interface.Context().Get(ctx).MerchantMember.Id)
	}
	update, err := service.SubscriptionUpdate(ctx, &subscription2.UpdateReq{
		SubscriptionId:      req.SubscriptionId,
		NewPlanId:           req.NewPlanId,
		Quantity:            req.Quantity,
		AddonParams:         req.AddonParams,
		WithImmediateEffect: req.WithImmediateEffect,
		ConfirmTotalAmount:  req.ConfirmTotalAmount,
		ConfirmCurrency:     req.ConfirmCurrency,
		ProrationDate:       req.ProrationDate,
		Metadata:            req.Metadata,
	}, merchantMemberId)
	if err != nil {
		return nil, err
	}
	return &subscription.UpdateRes{
		SubscriptionPendingUpdate: update.SubscriptionPendingUpdate,
		Paid:                      update.Paid,
		Link:                      update.Link,
	}, nil
}

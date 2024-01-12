package merchant

import (
	"context"
	subscription2 "go-oversea-pay/api/user/subscription"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq) (res *subscription.SubscriptionUpdateRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).Merchant.Id > 0, "merchantUserId invalid")
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
	}, int64(_interface.BizCtx().Get(ctx).Merchant.Id))
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdateRes{
		SubscriptionPendingUpdate: update.SubscriptionPendingUpdate,
		Paid:                      update.Paid,
		Link:                      update.Link,
	}, nil
}

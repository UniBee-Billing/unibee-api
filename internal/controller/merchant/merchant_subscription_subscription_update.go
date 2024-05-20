package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error) {
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		one := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx))
		utility.Assert(one != nil, "no active or incomplete subscription found")
		req.SubscriptionId = one.SubscriptionId
	}
	update, err := service.SubscriptionUpdate(ctx, &service.UpdateInternalReq{
		SubscriptionId:     req.SubscriptionId,
		NewPlanId:          req.NewPlanId,
		Quantity:           req.Quantity,
		AddonParams:        req.AddonParams,
		EffectImmediate:    req.EffectImmediate,
		GatewayId:          req.GatewayId,
		ConfirmTotalAmount: req.ConfirmTotalAmount,
		ConfirmCurrency:    req.ConfirmCurrency,
		ProrationDate:      req.ProrationDate,
		Metadata:           req.Metadata,
		DiscountCode:       req.DiscountCode,
		TaxPercentage:      req.TaxPercentage,
		Discount:           req.Discount,
		ManualPayment:      req.ManualPayment,
		ReturnUrl:          req.ReturnUrl,
		ProductData:        req.ProductData,
	}, -1)
	if err != nil {
		return nil, err
	}
	return &subscription.UpdateRes{
		SubscriptionPendingUpdate: update.SubscriptionPendingUpdate,
		Paid:                      update.Paid,
		Link:                      update.Link,
	}, nil
}

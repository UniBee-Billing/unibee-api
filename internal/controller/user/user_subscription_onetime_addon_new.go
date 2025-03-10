package user

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/onetime"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) OnetimeAddonNew(ctx context.Context, req *subscription.OnetimeAddonNewReq) (res *subscription.OnetimeAddonNewRes, err error) {
	result, err := onetime.CreateSubOneTimeAddon(ctx, &onetime.SubscriptionCreateOnetimeAddonInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		SubscriptionId: req.SubscriptionId,
		AddonId:        req.AddonId,
		Quantity:       req.Quantity,
		RedirectUrl:    req.ReturnUrl,
		Metadata:       req.Metadata,
		DiscountCode:   req.DiscountCode,
		GatewayId:      req.GatewayId,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.OnetimeAddonNewRes{
		SubscriptionOnetimeAddon: result.SubscriptionOnetimeAddon,
		Paid:                     result.Paid,
		Link:                     result.Link,
		Invoice:                  result.Invoice,
	}, nil
}

package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/onetime"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) OnetimeAddonNew(ctx context.Context, req *subscription.OnetimeAddonNewReq) (res *subscription.OnetimeAddonNewRes, err error) {
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		one := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx))
		utility.Assert(one != nil, "no active or incomplete subscription found")
		req.SubscriptionId = one.SubscriptionId
	}
	result, err := onetime.CreateSubOneTimeAddon(ctx, &onetime.SubscriptionCreateOnetimeAddonInternalReq{
		MerchantId:         _interface.GetMerchantId(ctx),
		SubscriptionId:     req.SubscriptionId,
		AddonId:            req.AddonId,
		Quantity:           req.Quantity,
		RedirectUrl:        req.ReturnUrl,
		Metadata:           req.Metadata,
		DiscountCode:       req.DiscountCode,
		DiscountAmount:     req.DiscountAmount,
		DiscountPercentage: req.DiscountPercentage,
		TaxPercentage:      req.TaxPercentage,
		GatewayId:          req.GatewayId,
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

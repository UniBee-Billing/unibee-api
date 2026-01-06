package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/onetime"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) OnetimeAddonPurchaseList(ctx context.Context, req *subscription.OnetimeAddonPurchaseListReq) (res *subscription.OnetimeAddonPurchaseListRes, err error) {
	return &subscription.OnetimeAddonPurchaseListRes{SubscriptionOnetimeAddons: onetime.SubscriptionOnetimeAddonPurchaseList(ctx, &onetime.SubscriptionOnetimeAddonPurchaseListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		UserId:     req.UserId,
		Page:       req.Page,
		Count:      req.Count,
	})}, nil
}

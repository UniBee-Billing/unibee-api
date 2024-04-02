package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/onetime"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) OnetimeAddonList(ctx context.Context, req *subscription.OnetimeAddonListReq) (res *subscription.OnetimeAddonListRes, err error) {
	return &subscription.OnetimeAddonListRes{SubscriptionOnetimeAddons: onetime.SubscriptionOnetimeAddonList(ctx, &onetime.SubscriptionOnetimeAddonListInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		SubscriptionId: req.SubscriptionId,
		Page:           req.Page,
		Count:          req.Count,
	})}, nil
}

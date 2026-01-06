package user

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/onetime"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) OnetimeAddonList(ctx context.Context, req *subscription.OnetimeAddonListReq) (res *subscription.OnetimeAddonListRes, err error) {
	return &subscription.OnetimeAddonListRes{SubscriptionOnetimeAddons: onetime.SubscriptionOnetimeAddonPurchaseList(ctx, &onetime.SubscriptionOnetimeAddonPurchaseListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		UserId:     _interface.Context().Get(ctx).User.Id,
		Page:       req.Page,
		Count:      req.Count,
	})}, nil
}

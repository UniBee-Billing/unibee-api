package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Renew(ctx context.Context, req *subscription.RenewReq) (res *subscription.RenewRes, err error) {
	renewRes, err := service.SubscriptionRenew(ctx, &service.RenewInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		SubscriptionId: req.SubscriptionId,
		UserId:         req.UserId,
		GatewayId:      req.GatewayId,
	})
	return &subscription.RenewRes{
		Subscription: renewRes.Subscription,
		Paid:         renewRes.Paid,
		Link:         renewRes.Link,
	}, nil
}

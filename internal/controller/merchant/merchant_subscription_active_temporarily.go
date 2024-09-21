package merchant

import (
	"context"
	"unibee/internal/logic/subscription/service"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) ActiveTemporarily(ctx context.Context, req *subscription.ActiveTemporarilyReq) (res *subscription.ActiveTemporarilyRes, err error) {
	err = service.SubscriptionActiveTemporarily(ctx, req.SubscriptionId, req.ExpireTime)
	if err != nil {
		return nil, err
	}
	return &subscription.ActiveTemporarilyRes{}, nil
}

package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) Cancel(ctx context.Context, req *subscription.CancelReq) (res *subscription.CancelRes, err error) {

	err = service.SubscriptionCancel(ctx, req.SubscriptionId, req.Prorate, req.InvoiceNow, "Admin Cancel")
	if err != nil {
		return nil, err
	}
	return &subscription.CancelRes{}, nil
}

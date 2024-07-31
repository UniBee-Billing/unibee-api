package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) Cancel(ctx context.Context, req *subscription.CancelReq) (res *subscription.CancelRes, err error) {
	if len(req.SubscriptionId) == 0 {
		utility.Assert(req.UserId > 0, "one of SubscriptionId and UserId should provide")
		one := query.GetLatestActiveOrIncompleteSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx), req.ProductId)
		utility.Assert(one != nil, "no active or incomplete subscription found")
		req.SubscriptionId = one.SubscriptionId
	}
	err = service.SubscriptionCancel(ctx, req.SubscriptionId, req.Prorate, req.InvoiceNow, "Admin Cancel")
	if err != nil {
		return nil, err
	}
	return &subscription.CancelRes{}, nil
}

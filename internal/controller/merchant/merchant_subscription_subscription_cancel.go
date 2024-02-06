package merchant

import (
	"context"
	"unibee-api/api/merchant/subscription"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/subscription/service"
	"unibee-api/utility"
)

func (c *ControllerSubscription) SubscriptionCancel(ctx context.Context, req *subscription.SubscriptionCancelReq) (res *subscription.SubscriptionCancelRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}

	err = service.SubscriptionCancel(ctx, req.SubscriptionId, req.Prorate, req.InvoiceNow, "Admin Cancel")
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionCancelRes{}, nil
}

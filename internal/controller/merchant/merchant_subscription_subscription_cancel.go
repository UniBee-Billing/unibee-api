package merchant

import (
	"context"
	"unibee/api/merchant/subscription"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"
)

func (c *ControllerSubscription) Cancel(ctx context.Context, req *subscription.CancelReq) (res *subscription.CancelRes, err error) {

	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}

	err = service.SubscriptionCancel(ctx, req.SubscriptionId, req.Prorate, req.InvoiceNow, "Admin Cancel")
	if err != nil {
		return nil, err
	}
	return &subscription.CancelRes{}, nil
}

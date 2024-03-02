package merchant

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionAddNewTrialStart(ctx context.Context, req *subscription.SubscriptionAddNewTrialStartReq) (res *subscription.SubscriptionAddNewTrialStartRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	err = service.SubscriptionAddNewTrialEnd(ctx, req.SubscriptionId, req.AppendTrialEndHour)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionAddNewTrialStartRes{}, nil
}

package system

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"

	"unibee/api/system/subscription"
)

func (c *ControllerSubscription) SubscriptionEndTrial(ctx context.Context, req *subscription.SubscriptionEndTrialReq) (res *subscription.SubscriptionEndTrialRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantMember.Id > 0, "merchantMemberId invalid")
	}
	err = service.SubscriptionEndTrial(ctx, req.SubscriptionId)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionEndTrialRes{}, nil
}

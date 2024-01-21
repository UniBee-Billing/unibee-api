package merchant

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/merchant/subscription"
)

func (c *ControllerSubscription) UserSubscriptionDetail(ctx context.Context, req *subscription.UserSubscriptionDetailReq) (res *subscription.UserSubscriptionDetailRes, err error) {
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser != nil, "merchant auth failure,not login")
		utility.Assert(_interface.BizCtx().Get(ctx).MerchantUser.Id > 0, "merchantUserId invalid")
	}
	user := query.GetUserAccountById(ctx, uint64(req.UserId))
	utility.Assert(user != nil, "user not found")
	if user != nil {
		user.Password = ""
	}
	one := query.GetLatestSubscriptionByUserId(ctx, req.UserId, req.MerchantId)
	if one != nil {
		detail, err := service.SubscriptionDetail(ctx, one.SubscriptionId)
		if err == nil {
			return &subscription.UserSubscriptionDetailRes{
				User:                                user,
				Subscription:                        detail.Subscription,
				Plan:                                detail.Plan,
				Channel:                             detail.Channel,
				Addons:                              detail.Addons,
				UnfinishedSubscriptionPendingUpdate: detail.UnfinishedSubscriptionPendingUpdate,
			}, nil
		}
	}

	return &subscription.UserSubscriptionDetailRes{
		User:                                user,
		Subscription:                        nil,
		Plan:                                nil,
		Channel:                             nil,
		Addons:                              nil,
		UnfinishedSubscriptionPendingUpdate: nil,
	}, nil
}

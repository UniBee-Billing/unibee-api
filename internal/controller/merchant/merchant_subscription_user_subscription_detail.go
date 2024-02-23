package merchant

import (
	"context"
	"unibee-api/internal/consts"
	_interface "unibee-api/internal/interface"
	"unibee-api/internal/logic/gateway/ro"
	"unibee-api/internal/logic/subscription/service"
	"unibee-api/internal/query"
	"unibee-api/utility"

	"unibee-api/api/merchant/subscription"
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
	one := query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, req.UserId, _interface.GetMerchantId(ctx))
	if one != nil {
		detail, err := service.SubscriptionDetail(ctx, one.SubscriptionId)
		if err == nil {
			return &subscription.UserSubscriptionDetailRes{
				User:                                detail.User,
				Subscription:                        detail.Subscription,
				Plan:                                detail.Plan,
				Gateway:                             detail.Gateway,
				Addons:                              detail.Addons,
				UnfinishedSubscriptionPendingUpdate: detail.UnfinishedSubscriptionPendingUpdate,
			}, nil
		}
	}

	return &subscription.UserSubscriptionDetailRes{
		User:                                ro.SimplifyUserAccount(user),
		Subscription:                        nil,
		Plan:                                nil,
		Gateway:                             nil,
		Addons:                              nil,
		UnfinishedSubscriptionPendingUpdate: nil,
	}, nil
}

package merchant

import (
	"context"
	"unibee/api/bean"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) UserSubscriptionDetail(ctx context.Context, req *subscription.UserSubscriptionDetailReq) (res *subscription.UserSubscriptionDetailRes, err error) {
	var user *entity.UserAccount
	if _interface.Context().Get(ctx).IsOpenApiCall {
		if req.UserId == 0 {
			utility.Assert(len(req.ExternalUserId) > 0, "ExternalUserId|UserId is nil, one of it is required")
			user = query.GetUserAccountByExternalUserId(ctx, _interface.GetMerchantId(ctx), req.ExternalUserId)
			utility.AssertError(err, "Server Error")
		} else {
			user = query.GetUserAccountById(ctx, req.UserId)
		}
	} else {
		user = query.GetUserAccountById(ctx, req.UserId)
	}
	utility.Assert(user != nil, "user not found")
	one := query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, user.Id, _interface.GetMerchantId(ctx), req.ProductId)
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
		User:                                bean.SimplifyUserAccount(user),
		Subscription:                        nil,
		Plan:                                nil,
		Gateway:                             nil,
		Addons:                              nil,
		UnfinishedSubscriptionPendingUpdate: nil,
	}, nil
}

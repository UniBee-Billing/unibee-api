package user

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) UserCurrentSubscriptionDetail(ctx context.Context, req *subscription.UserCurrentSubscriptionDetailReq) (res *subscription.UserCurrentSubscriptionDetailRes, err error) {
	user := query.GetUserAccountById(ctx, _interface.Context().Get(ctx).User.Id)
	one := query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, user.Id, _interface.GetMerchantId(ctx), req.ProductId)
	if one != nil {
		detail, err := service.SubscriptionDetail(ctx, one.SubscriptionId)
		if err != nil {
			return nil, err
		}
		if detail != nil {
			return &subscription.UserCurrentSubscriptionDetailRes{
				User:                                detail.User,
				Subscription:                        detail.Subscription,
				Plan:                                detail.Plan,
				Gateway:                             detail.Gateway,
				AddonParams:                         detail.AddonParams,
				Addons:                              detail.Addons,
				LatestInvoice:                       detail.LatestInvoice,
				UnfinishedSubscriptionPendingUpdate: detail.UnfinishedSubscriptionPendingUpdate,
			}, nil
		} else {
			return nil, nil
		}
	}
	return nil, nil
}

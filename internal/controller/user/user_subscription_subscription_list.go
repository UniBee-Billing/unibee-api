package user

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface/context"
	detail2 "unibee/internal/logic/subscription/service/detail"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) List(ctx context.Context, req *subscription.ListReq) (res *subscription.ListRes, err error) {
	// return one latest user subscription list as unique subscription
	var subDetails []*detail.SubscriptionDetail
	subs := query.GetLatestActiveOrIncompleteOrCreateSubscriptionsByUserId(ctx, _interface.Context().Get(ctx).User.Id, _interface.GetMerchantId(ctx))
	for _, sub := range subs {
		if sub != nil {
			subDetailRes, err := detail2.SubscriptionDetail(ctx, sub.SubscriptionId)
			if err == nil {
				var addonParams []*bean.PlanAddonParam
				_ = utility.UnmarshalFromJsonString(sub.AddonData, &addonParams)
				subDetails = append(subDetails, &detail.SubscriptionDetail{
					User:                                subDetailRes.User,
					Subscription:                        subDetailRes.Subscription,
					Plan:                                subDetailRes.Plan,
					Gateway:                             subDetailRes.Gateway,
					AddonParams:                         addonParams,
					Addons:                              subDetailRes.Addons,
					LatestInvoice:                       subDetailRes.LatestInvoice,
					Discount:                            subDetailRes.Discount,
					UnfinishedSubscriptionPendingUpdate: subDetailRes.UnfinishedSubscriptionPendingUpdate,
				})
			}
		}
	}
	return &subscription.ListRes{Subscriptions: subDetails}, nil
}

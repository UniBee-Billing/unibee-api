package user

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/subscription"
)

func (c *ControllerSubscription) List(ctx context.Context, req *subscription.ListReq) (res *subscription.ListRes, err error) {
	// return one latest user subscription list as unique subscription
	var subDetails []*detail.SubscriptionDetail
	sub := query.GetLatestActiveOrIncompleteOrCreateSubscriptionByUserId(ctx, int64(_interface.BizCtx().Get(ctx).User.Id), _interface.GetMerchantId(ctx))
	if sub != nil {
		subDetailRes, err := service.SubscriptionDetail(ctx, sub.SubscriptionId)
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
				UnfinishedSubscriptionPendingUpdate: subDetailRes.UnfinishedSubscriptionPendingUpdate,
			})
		}
	}
	return &subscription.ListRes{Subscriptions: subDetails}, nil
}

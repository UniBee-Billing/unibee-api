package user

import (
	"context"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/gateway/ro"
	"unibee/internal/logic/subscription/service"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/user/subscription"
)

// SubscriptionList todo mark demo requirement, return only one user sub by gmt_create desc
func (c *ControllerSubscription) List(ctx context.Context, req *subscription.ListReq) (res *subscription.ListRes, err error) {
	// service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == req.UserId, "userId not match")
	}
	// return one latest user subscription list as unique subscription
	var subDetails []*ro.SubscriptionDetailVo
	sub := query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, int64(_interface.BizCtx().Get(ctx).User.Id), _interface.GetMerchantId(ctx))
	if sub != nil {
		subDetailRes, err := service.SubscriptionDetail(ctx, sub.SubscriptionId)
		if err == nil {
			var addonParams []*ro.SubscriptionPlanAddonParamRo
			_ = utility.UnmarshalFromJsonString(sub.AddonData, &addonParams)
			subDetails = append(subDetails, &ro.SubscriptionDetailVo{
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

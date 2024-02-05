package user

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/gateway/ro"
	"go-oversea-pay/internal/logic/subscription/service"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/user/subscription"
)

// SubscriptionList todo mark demo requirement, return only one user sub by gmt_create desc
func (c *ControllerSubscription) SubscriptionList(ctx context.Context, req *subscription.SubscriptionListReq) (res *subscription.SubscriptionListRes, err error) {
	// service 层不做用户校验
	if !consts.GetConfigInstance().IsLocal() {
		//User 检查
		utility.Assert(_interface.BizCtx().Get(ctx).User != nil, "auth failure,not login")
		utility.Assert(int64(_interface.BizCtx().Get(ctx).User.Id) == req.UserId, "userId not match")
	}
	// return one latest user subscription list as unique subscription
	var subDetails []*ro.SubscriptionDetailRo
	sub := query.GetLatestActiveOrCreateSubscriptionByUserId(ctx, int64(_interface.BizCtx().Get(ctx).User.Id), req.MerchantId)
	if sub != nil {
		subDetailRes, err := service.SubscriptionDetail(ctx, sub.SubscriptionId)
		if err == nil {
			var addonParams []*ro.SubscriptionPlanAddonParamRo
			_ = utility.UnmarshalFromJsonString(sub.AddonData, &addonParams)
			subDetails = append(subDetails, &ro.SubscriptionDetailRo{
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
	return &subscription.SubscriptionListRes{Subscriptions: subDetails}, nil
	//return &subscription.SubscriptionListRes{Subscriptions: service.SubscriptionList(ctx, &service.SubscriptionListInternalReq{
	//	MerchantId: req.MerchantId,
	//	UserId:     req.UserId,
	//	Status:     consts.SubStatusActive,
	//	SortField:  "gmt_create",
	//	SortType:   "desc",
	//	Page:       0,
	//	Count:      1,
	//})}, nil
}

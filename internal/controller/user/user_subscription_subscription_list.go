package user

import (
	"context"
	"go-oversea-pay/internal/consts"
	_interface "go-oversea-pay/internal/interface"
	"go-oversea-pay/internal/logic/subscription/service"
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
	return &subscription.SubscriptionListRes{Subscriptions: service.SubscriptionList(ctx, &service.SubscriptionListInternalReq{
		MerchantId: req.MerchantId,
		UserId:     req.UserId,
		Status:     consts.SubStatusActive,
		SortField:  "gmt_create",
		SortType:   "desc",
		Page:       0,
		Count:      1,
	})}, nil
}

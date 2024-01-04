package user

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionList(ctx context.Context, req *subscription.SubscriptionListReq) (res *subscription.SubscriptionListRes, err error) {
	return &subscription.SubscriptionListRes{Subscriptions: service.SubscriptionList(ctx, &service.SubscriptionListInternalReq{
		MerchantId: req.MerchantId,
		UserId:     req.UserId,
		Status:     req.Status,
		Page:       req.Page,
		Count:      req.Count,
	})}, nil
}

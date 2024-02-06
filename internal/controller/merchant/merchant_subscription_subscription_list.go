package merchant

import (
	"context"
	"unibee-api/internal/logic/subscription/service"

	"unibee-api/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionList(ctx context.Context, req *subscription.SubscriptionListReq) (res *subscription.SubscriptionListRes, err error) {
	return &subscription.SubscriptionListRes{Subscriptions: service.SubscriptionList(ctx, &service.SubscriptionListInternalReq{
		MerchantId: req.MerchantId,
		UserId:     req.UserId,
		Status:     req.Status,
		SortField:  req.SortField,
		SortType:   req.SortType,
		Page:       req.Page,
		Count:      req.Count,
	})}, nil
}

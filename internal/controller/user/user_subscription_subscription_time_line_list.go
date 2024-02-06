package user

import (
	"context"
	"unibee-api/api/user/subscription"
	"unibee-api/internal/logic/subscription/service"
)

func (c *ControllerSubscription) SubscriptionTimeLineList(ctx context.Context, req *subscription.SubscriptionTimeLineListReq) (res *subscription.SubscriptionTimeLineListRes, err error) {
	result, err := service.SubscriptionTimeLineList(ctx, &service.SubscriptionTimeLineListInternalReq{
		MerchantId: req.MerchantId,
		UserId:     req.UserId,
		SortField:  req.SortField,
		SortType:   req.SortType,
		Page:       req.Page,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionTimeLineListRes{SubscriptionTimeLines: result.SubscriptionTimelines}, nil
}

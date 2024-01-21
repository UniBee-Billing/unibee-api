package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/merchant/subscription"
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

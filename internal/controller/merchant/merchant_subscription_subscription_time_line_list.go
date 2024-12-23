package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/timeline"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) TimeLineList(ctx context.Context, req *subscription.TimeLineListReq) (res *subscription.TimeLineListRes, err error) {
	result, err := timeline.SubscriptionTimeLineList(ctx, &timeline.SubscriptionTimeLineListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		UserId:     req.UserId,
		SortField:  req.SortField,
		SortType:   req.SortType,
		Page:       req.Page,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.TimeLineListRes{SubscriptionTimeLines: result.SubscriptionTimelines, Total: result.Total}, nil
}

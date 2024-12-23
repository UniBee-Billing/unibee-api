package user

import (
	"context"
	"unibee/api/user/subscription"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/timeline"
)

func (c *ControllerSubscription) TimeLineList(ctx context.Context, req *subscription.TimeLineListReq) (res *subscription.TimeLineListRes, err error) {
	result, err := timeline.SubscriptionTimeLineList(ctx, &timeline.SubscriptionTimeLineListInternalReq{
		MerchantId: _interface.GetMerchantId(ctx),
		UserId:     _interface.Context().Get(ctx).User.Id,
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

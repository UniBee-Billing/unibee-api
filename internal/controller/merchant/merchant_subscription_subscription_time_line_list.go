package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) TimeLineList(ctx context.Context, req *subscription.TimeLineListReq) (res *subscription.TimeLineListRes, err error) {
	result, err := service.SubscriptionTimeLineList(ctx, &service.SubscriptionTimeLineListInternalReq{
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

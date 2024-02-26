package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionMerchantPendingUpdateList(ctx context.Context, req *subscription.SubscriptionMerchantPendingUpdateListReq) (res *subscription.SubscriptionMerchantPendingUpdateListRes, err error) {
	result, err := service.SubscriptionPendingUpdateList(ctx, &service.SubscriptionPendingUpdateListInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		SubscriptionId: req.SubscriptionId,
		SortField:      req.SortField,
		SortType:       req.SortType,
		Page:           req.Page,
		Count:          req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionMerchantPendingUpdateListRes{SubscriptionPendingUpdateDetails: result.SubscriptionPendingUpdateDetails}, nil
}

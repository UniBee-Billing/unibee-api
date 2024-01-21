package merchant

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/merchant/subscription"
)

func (c *ControllerSubscription) SubscriptionMerchantPendingUpdateList(ctx context.Context, req *subscription.SubscriptionMerchantPendingUpdateListReq) (res *subscription.SubscriptionMerchantPendingUpdateListRes, err error) {
	result, err := service.SubscriptionPendingUpdateList(ctx, &service.SubscriptionPendingUpdateListInternalReq{
		MerchantId:     req.MerchantId,
		SubscriptionId: req.SubscriptionId,
		SortField:      req.SortField,
		SortType:       req.SortType,
		Page:           req.Page,
		Count:          req.Count,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionMerchantPendingUpdateListRes{SubscriptionPendingUpdates: result.SubscriptionPendingUpdates}, nil
}

package merchant

import (
	"context"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/service"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) List(ctx context.Context, req *subscription.ListReq) (res *subscription.ListRes, err error) {
	list, total := service.SubscriptionList(ctx, &service.SubscriptionListInternalReq{
		MerchantId:  _interface.GetMerchantId(ctx),
		UserId:      req.UserId,
		Status:      req.Status,
		Currency:    req.Currency,
		AmountStart: req.AmountStart,
		AmountEnd:   req.AmountEnd,
		PlanIds:     req.PlanIds,
		ProductIds:  req.ProductIds,
		SortField:   req.SortField,
		SortType:    req.SortType,
		Page:        req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
	})
	return &subscription.ListRes{Subscriptions: list, Total: total}, nil
}

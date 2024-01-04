package user

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdatePrepare(ctx context.Context, req *subscription.SubscriptionUpdatePrepareReq) (res *subscription.SubscriptionUpdatePrepareRes, err error) {
	prepare, err := service.SubscriptionUpdatePreview(ctx, req)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdatePrepareRes{
		TotalAmount:   prepare.TotalAmount,
		Currency:      prepare.Currency,
		Invoice:       prepare.Invoice,
		ProrationDate: prepare.ProrationDate,
	}, nil
}

package user

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq) (res *subscription.SubscriptionUpdatePreviewRes, err error) {
	prepare, err := service.SubscriptionUpdatePreview(ctx, req)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionUpdatePreviewRes{
		TotalAmount:   prepare.TotalAmount,
		Currency:      prepare.Currency,
		Invoice:       prepare.Invoice,
		ProrationDate: prepare.ProrationDate,
	}, nil
}

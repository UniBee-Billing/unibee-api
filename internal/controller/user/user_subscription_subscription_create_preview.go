package user

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionCreatePreview(ctx context.Context, req *subscription.SubscriptionCreatePreviewReq) (res *subscription.SubscriptionCreatePreviewRes, err error) {
	prepare, err := service.SubscriptionCreatePreview(ctx, req)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionCreatePreviewRes{
		Plan:        prepare.Plan,
		Quantity:    prepare.Quantity,
		PayChannel:  prepare.PayChannel,
		AddonParams: prepare.AddonParams,
		Addons:      prepare.Addons,
		TotalAmount: prepare.TotalAmount,
		Currency:    prepare.Currency,
		Invoice:     prepare.Invoice,
		UserId:      prepare.UserId,
		Email:       prepare.Email,
	}, nil
}

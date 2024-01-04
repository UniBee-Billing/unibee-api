package user

import (
	"context"
	"go-oversea-pay/internal/logic/subscription/service"

	"go-oversea-pay/api/user/subscription"
)

func (c *ControllerSubscription) SubscriptionCreatePrepare(ctx context.Context, req *subscription.SubscriptionCreatePrepareReq) (res *subscription.SubscriptionCreatePrepareRes, err error) {
	prepare, err := service.SubscriptionCreatePrepare(ctx, req)
	if err != nil {
		return nil, err
	}
	return &subscription.SubscriptionCreatePrepareRes{
		Plan:        prepare.Plan,
		Quantity:    prepare.Quantity,
		PayChannel:  prepare.PayChannel,
		AddonParams: prepare.AddonParams,
		Addons:      prepare.Addons,
		TotalAmount: prepare.TotalAmount,
		Currency:    prepare.Currency,
		UserId:      prepare.UserId,
		Email:       prepare.Email,
	}, nil
}

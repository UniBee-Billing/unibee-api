package merchant

import (
	"context"
	"unibee/internal/logic/subscription/service"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error) {
	createRes, err := service.SubscriptionCreate(ctx, &service.CreateInternalReq{
		PlanId:             req.PlanId,
		UserId:             req.UserId,
		Quantity:           req.Quantity,
		GatewayId:          req.GatewayId,
		AddonParams:        req.AddonParams,
		ConfirmTotalAmount: req.ConfirmTotalAmount,
		ConfirmCurrency:    req.ConfirmCurrency,
		ReturnUrl:          req.ReturnUrl,
		VatCountryCode:     req.VatCountryCode,
		VatNumber:          req.VatNumber,
		PaymentMethodId:    req.PaymentMethodId,
		Metadata:           req.Metadata,
	})
	return &subscription.CreateRes{
		Subscription: createRes.Subscription,
		Paid:         createRes.Paid,
		Link:         createRes.Link,
	}, nil
}

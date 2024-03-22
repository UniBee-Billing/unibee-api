package user

import (
	"context"
	"unibee/api/user/subscription"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error) {
	createRes, err := service.SubscriptionCreate(ctx, &service.CreateInternalReq{
		PlanId:             req.PlanId,
		UserId:             _interface.BizCtx().Get(ctx).User.Id,
		MerchantId:         _interface.GetMerchantId(ctx),
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

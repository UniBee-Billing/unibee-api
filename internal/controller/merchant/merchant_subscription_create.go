package merchant

import (
	"context"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/auth"
	"unibee/internal/logic/subscription/service"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error) {
	if req.UserId == 0 {
		utility.Assert(len(req.Email) > 0, "Email|UserId is nil")
		user, err := auth.QueryOrCreateUser(ctx, &auth.NewReq{
			ExternalUserId: req.ExternalUserId,
			Email:          req.Email,
			MerchantId:     _interface.GetMerchantId(ctx),
		})
		utility.AssertError(err, "Server Error")
		req.UserId = user.Id
	}
	utility.Assert(req.UserId > 0, "Invalid UserId")
	createRes, err := service.SubscriptionCreate(ctx, &service.CreateInternalReq{
		MerchantId:         _interface.GetMerchantId(ctx),
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
		TaxPercentage:      req.TaxPercentage,
		PaymentMethodId:    req.PaymentMethodId,
		Metadata:           req.Metadata,
		DiscountCode:       req.DiscountCode,
		Discount:           req.Discount,
		TrialEnd:           req.TrialEnd,
	})
	return &subscription.CreateRes{
		Subscription: createRes.Subscription,
		Paid:         createRes.Paid,
		Link:         createRes.Link,
	}, nil
}

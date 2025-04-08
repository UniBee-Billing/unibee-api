package checkout

import (
	"context"
	subscription2 "unibee/api/merchant/subscription"
	"unibee/internal/controller/merchant"
	_interface "unibee/internal/interface/context"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/checkout/subscription"
)

func (c *ControllerSubscription) Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error) {
	utility.Assert(req.PlanId > 0, "PlanId is required")
	plan := query.GetPlanById(ctx, req.PlanId)
	utility.Assert(plan != nil, "Plan not found")
	_interface.Context().Get(ctx).MerchantId = plan.MerchantId
	controllerSubscription := merchant.ControllerSubscription{}
	createRes, err := controllerSubscription.Create(ctx, &subscription2.CreateReq{
		PlanId:                 req.PlanId,
		Email:                  req.Email,
		UserId:                 req.UserId,
		ExternalUserId:         req.ExternalUserId,
		User:                   req.User,
		Quantity:               req.Quantity,
		GatewayId:              req.GatewayId,
		GatewayPaymentType:     req.GatewayPaymentType,
		AddonParams:            req.AddonParams,
		ReturnUrl:              req.ReturnUrl,
		CancelUrl:              req.CancelUrl,
		VatCountryCode:         req.VatCountryCode,
		VatNumber:              req.VatNumber,
		PaymentMethodId:        req.PaymentMethodId,
		TaxPercentage:          req.TaxPercentage,
		Metadata:               req.Metadata,
		DiscountCode:           req.DiscountCode,
		Discount:               req.Discount,
		TrialEnd:               req.TrialEnd,
		StartIncomplete:        req.StartIncomplete,
		ProductData:            req.ProductData,
		ApplyPromoCredit:       req.ApplyPromoCredit,
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
		ConfirmCurrency:        req.ConfirmCurrency,
		ConfirmTotalAmount:     req.ConfirmTotalAmount,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.CreateRes{
		Subscription:                   createRes.Subscription,
		User:                           createRes.User,
		Paid:                           createRes.Paid,
		Link:                           createRes.Link,
		Token:                          createRes.Token,
		OtherPendingCryptoSubscription: createRes.OtherPendingCryptoSubscription,
	}, nil
}

package merchant

import (
	"context"
	"fmt"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/subscription/service"
	user2 "unibee/internal/logic/user"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error) {
	if req.UserId == 0 && req.User != nil {
		user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
			ExternalUserId: req.User.ExternalUserId,
			Email:          req.User.Email,
			FirstName:      req.User.FirstName,
			LastName:       req.User.LastName,
			Phone:          req.User.Phone,
			Address:        req.User.Address,
			UserName:       req.User.UserName,
			CountryCode:    req.User.CountryCode,
			Type:           req.User.Type,
			CompanyName:    req.User.CompanyName,
			VATNumber:      req.User.VatNumber,
			City:           req.User.City,
			ZipCode:        req.User.ZipCode,
			MerchantId:     _interface.GetMerchantId(ctx),
		})
		utility.AssertError(err, "Server Error")
		req.UserId = user.Id
	} else if req.UserId == 0 && len(req.Email) > 0 {
		user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
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
		CancelUrl:          req.CancelUrl,
		VatCountryCode:     req.VatCountryCode,
		VatNumber:          req.VatNumber,
		TaxPercentage:      req.TaxPercentage,
		PaymentMethodId:    req.PaymentMethodId,
		Metadata:           req.Metadata,
		DiscountCode:       req.DiscountCode,
		Discount:           req.Discount,
		TrialEnd:           req.TrialEnd,
		StartIncomplete:    req.StartIncomplete,
		ProductData:        req.ProductData,
	})
	if err == nil {
		operation_log.AppendOptLog(ctx, &operation_log.OptLogRequest{
			MerchantId:     createRes.Subscription.MerchantId,
			Target:         fmt.Sprintf("Subscription(%v)", createRes.Subscription.SubscriptionId),
			Content:        "AssignSubscriptionToUser",
			UserId:         createRes.Subscription.UserId,
			SubscriptionId: createRes.Subscription.SubscriptionId,
			InvoiceId:      "",
			PlanId:         0,
			DiscountCode:   "",
		}, err)
	}

	return &subscription.CreateRes{
		Subscription: createRes.Subscription,
		Paid:         createRes.Paid,
		Link:         createRes.Link,
	}, nil
}

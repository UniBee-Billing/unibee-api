package merchant

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	user2 "unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error) {
	if req.UserId == 0 && req.User != nil {
		user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
			ExternalUserId:     req.User.ExternalUserId,
			Email:              req.User.Email,
			FirstName:          req.User.FirstName,
			LastName:           req.User.LastName,
			Phone:              req.User.Phone,
			Address:            req.User.Address,
			UserName:           req.User.UserName,
			CountryCode:        req.User.CountryCode,
			Type:               req.User.Type,
			CompanyName:        req.User.CompanyName,
			VATNumber:          req.User.VatNumber,
			City:               req.User.City,
			ZipCode:            req.User.ZipCode,
			Language:           req.User.Language,
			RegistrationNumber: req.User.RegistrationNumber,
			MerchantId:         _interface.GetMerchantId(ctx),
		})
		utility.AssertError(err, "Server Error")
		req.UserId = user.Id
	} else if req.UserId == 0 && len(req.Email) > 0 {
		utility.Assert(len(req.Email) > 0, "Email|UserId is nil")
		user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
			ExternalUserId: req.ExternalUserId,
			Email:          req.Email,
			MerchantId:     _interface.GetMerchantId(ctx),
		})
		utility.AssertError(err, "Server Error")
		req.UserId = user.Id
	}
	prepare, err := service.SubscriptionCreatePreview(ctx, &service.CreatePreviewInternalReq{
		MerchantId:     _interface.GetMerchantId(ctx),
		PlanId:         req.PlanId,
		UserId:         req.UserId,
		Quantity:       req.Quantity,
		GatewayId:      req.GatewayId,
		AddonParams:    req.AddonParams,
		VatCountryCode: req.VatCountryCode,
		VatNumber:      req.VatNumber,
		TaxPercentage:  req.TaxPercentage,
		DiscountCode:   req.DiscountCode,
		TrialEnd:       req.TrialEnd,
	})
	if err != nil {
		return nil, err
	}
	var pendingCryptoSub *detail.SubscriptionDetail
	if prepare.UserId > 0 {
		one := query.GetLatestCreateOrProcessingSubscriptionByUserId(ctx, prepare.UserId, _interface.GetMerchantId(ctx), prepare.Plan.ProductId)
		if one != nil {
			gateway := query.GetGatewayById(ctx, one.GatewayId)
			if gateway.GatewayType == consts.GatewayTypeCrypto {
				pendingCryptoSub, _ = service.SubscriptionDetail(ctx, one.SubscriptionId)
			}
		}
	}
	return &subscription.CreatePreviewRes{
		Plan:              bean.SimplifyPlan(prepare.Plan),
		TrialEnd:          prepare.TrialEnd,
		Quantity:          prepare.Quantity,
		Gateway:           bean.SimplifyGateway(prepare.Gateway),
		AddonParams:       prepare.AddonParams,
		Addons:            prepare.Addons,
		OriginAmount:      prepare.OriginAmount,
		TotalAmount:       prepare.TotalAmount,
		DiscountAmount:    prepare.DiscountAmount,
		Currency:          prepare.Currency,
		VatNumber:         prepare.VatNumber,
		VatNumberValidate: prepare.VatNumberValidate,
		VatCountryCode:    prepare.VatCountryCode,
		VatCountryName:    prepare.VatCountryName,
		TaxPercentage:     prepare.TaxPercentage,
		Invoice: &bean.Invoice{
			InvoiceName:                    prepare.Invoice.InvoiceName,
			OriginAmount:                   prepare.Invoice.OriginAmount,
			TotalAmount:                    prepare.Invoice.TotalAmount,
			TotalAmountExcludingTax:        prepare.Invoice.TotalAmountExcludingTax,
			DiscountAmount:                 prepare.Invoice.DiscountAmount,
			DiscountCode:                   prepare.Invoice.DiscountCode,
			Currency:                       prepare.Invoice.Currency,
			TaxAmount:                      prepare.Invoice.TaxAmount,
			SubscriptionAmount:             prepare.Invoice.SubscriptionAmount,
			SubscriptionAmountExcludingTax: prepare.Invoice.SubscriptionAmountExcludingTax,
			Lines:                          prepare.Invoice.Lines,
		},
		UserId:                   prepare.UserId,
		Email:                    prepare.Email,
		Discount:                 prepare.Discount,
		VatNumberValidateMessage: prepare.VatNumberValidateMessage,
		DiscountMessage:          prepare.DiscountMessage,
		OtherPendingCryptoSub:    pendingCryptoSub,
	}, nil
}

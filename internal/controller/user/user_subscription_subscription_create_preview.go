package user

import (
	"context"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/user/subscription"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/subscription/service"
)

func (c *ControllerSubscription) CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error) {
	prepare, err := service.SubscriptionCreatePreview(ctx, &service.CreatePreviewInternalReq{
		MerchantId:             _interface.GetMerchantId(ctx),
		PlanId:                 req.PlanId,
		UserId:                 _interface.Context().Get(ctx).User.Id,
		Quantity:               req.Quantity,
		GatewayId:              req.GatewayId,
		AddonParams:            req.AddonParams,
		VatCountryCode:         req.VatCountryCode,
		VatNumber:              req.VatNumber,
		DiscountCode:           req.DiscountCode,
		ApplyPromoCredit:       req.ApplyPromoCredit,
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
	})
	if err != nil {
		return nil, err
	}
	return &subscription.CreatePreviewRes{
		Plan:                      bean.SimplifyPlan(prepare.Plan),
		TrialEnd:                  prepare.TrialEnd,
		Quantity:                  prepare.Quantity,
		Gateway:                   detail.ConvertGatewayDetail(ctx, prepare.Gateway),
		AddonParams:               prepare.AddonParams,
		Addons:                    prepare.Addons,
		OriginAmount:              prepare.OriginAmount,
		TotalAmount:               prepare.TotalAmount,
		DiscountAmount:            prepare.DiscountAmount,
		Currency:                  prepare.Currency,
		VatNumber:                 prepare.VatNumber,
		VatNumberValidate:         prepare.VatNumberValidate,
		VatCountryCode:            prepare.VatCountryCode,
		VatCountryName:            prepare.VatCountryName,
		TaxPercentage:             prepare.TaxPercentage,
		Invoice:                   prepare.Invoice,
		UserId:                    prepare.UserId,
		Email:                     prepare.Email,
		Discount:                  prepare.Discount,
		VatNumberValidateMessage:  prepare.VatNumberValidateMessage,
		DiscountMessage:           prepare.DiscountMessage,
		OtherActiveSubscriptionId: prepare.OtherActiveSubscriptionId,
		ApplyPromoCredit:          prepare.ApplyPromoCredit,
	}, nil
}

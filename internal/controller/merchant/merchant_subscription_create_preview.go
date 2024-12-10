package merchant

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean"
	"unibee/api/bean/detail"
	"unibee/api/merchant/subscription"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/subscription/service"
	user2 "unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/utility"
)

func (c *ControllerSubscription) CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error) {
	if req.UserId > 0 {
		user := query.GetUserAccountById(ctx, req.UserId)
		utility.Assert(user != nil, "user not found")
		if len(req.Email) > 0 {
			utility.Assert(user.Email == req.Email, "invalid email, not match")
		} else if req.User != nil {
			utility.Assert(user.Email == req.User.Email, "invalid user, not match")
			// Update user profile
			_, _ = dao.UserAccount.Ctx(ctx).Data(g.Map{
				dao.UserAccount.Columns().Address:            req.User.Address,
				dao.UserAccount.Columns().Phone:              req.User.Phone,
				dao.UserAccount.Columns().FirstName:          req.User.FirstName,
				dao.UserAccount.Columns().LastName:           req.User.LastName,
				dao.UserAccount.Columns().City:               req.User.City,
				dao.UserAccount.Columns().Type:               req.User.Type,
				dao.UserAccount.Columns().ZipCode:            req.User.ZipCode,
				dao.UserAccount.Columns().Language:           req.User.Language,
				dao.UserAccount.Columns().Address:            req.User.Address,
				dao.UserAccount.Columns().CompanyName:        req.User.CompanyName,
				dao.UserAccount.Columns().RegistrationNumber: req.User.RegistrationNumber,
				dao.UserAccount.Columns().GmtModify:          gtime.Now(),
			}).Where(dao.UserAccount.Columns().Id, req.UserId).OmitEmpty().Update()
		}
	} else if req.UserId == 0 && req.User != nil {
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
		MerchantId:             _interface.GetMerchantId(ctx),
		PlanId:                 req.PlanId,
		UserId:                 req.UserId,
		Quantity:               req.Quantity,
		GatewayId:              req.GatewayId,
		AddonParams:            req.AddonParams,
		VatCountryCode:         req.VatCountryCode,
		VatNumber:              req.VatNumber,
		TaxPercentage:          req.TaxPercentage,
		DiscountCode:           req.DiscountCode,
		TrialEnd:               req.TrialEnd,
		ApplyPromoCredit:       req.ApplyPromoCredit,
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
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
		Plan:                           bean.SimplifyPlan(prepare.Plan),
		TrialEnd:                       prepare.TrialEnd,
		Quantity:                       prepare.Quantity,
		Gateway:                        bean.SimplifyGateway(prepare.Gateway),
		AddonParams:                    prepare.AddonParams,
		Addons:                         prepare.Addons,
		TaxPercentage:                  prepare.TaxPercentage,
		SubscriptionAmountExcludingTax: prepare.Invoice.SubscriptionAmountExcludingTax,
		TaxAmount:                      prepare.Invoice.TaxAmount,
		DiscountAmount:                 prepare.DiscountAmount,
		TotalAmount:                    prepare.TotalAmount,
		OriginAmount:                   prepare.OriginAmount,
		Currency:                       prepare.Currency,
		VatNumber:                      prepare.VatNumber,
		VatNumberValidate:              prepare.VatNumberValidate,
		VatCountryCode:                 prepare.VatCountryCode,
		VatCountryName:                 prepare.VatCountryName,
		Invoice:                        prepare.Invoice,
		UserId:                         prepare.UserId,
		Email:                          prepare.Email,
		Discount:                       prepare.Discount,
		VatNumberValidateMessage:       prepare.VatNumberValidateMessage,
		DiscountMessage:                prepare.DiscountMessage,
		OtherPendingCryptoSubscription: pendingCryptoSub,
		OtherActiveSubscriptionId:      prepare.OtherActiveSubscriptionId,
		ApplyPromoCredit:               prepare.ApplyPromoCredit,
	}, nil
}

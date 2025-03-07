package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	dao "unibee/internal/dao/default"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/subscription/service"
	detail2 "unibee/internal/logic/subscription/service/detail"
	user2 "unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error) {
	token := ""
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
		MerchantId:             _interface.GetMerchantId(ctx),
		PlanId:                 req.PlanId,
		UserId:                 req.UserId,
		Quantity:               req.Quantity,
		GatewayId:              req.GatewayId,
		GatewayPaymentType:     req.GatewayPaymentType,
		AddonParams:            req.AddonParams,
		ConfirmTotalAmount:     req.ConfirmTotalAmount,
		ConfirmCurrency:        req.ConfirmCurrency,
		ReturnUrl:              req.ReturnUrl,
		CancelUrl:              req.CancelUrl,
		VatCountryCode:         req.VatCountryCode,
		VatNumber:              req.VatNumber,
		TaxPercentage:          req.TaxPercentage,
		PaymentMethodId:        req.PaymentMethodId,
		Metadata:               req.Metadata,
		DiscountCode:           req.DiscountCode,
		Discount:               req.Discount,
		TrialEnd:               req.TrialEnd,
		StartIncomplete:        req.StartIncomplete,
		ProductData:            req.ProductData,
		ApplyPromoCredit:       req.ApplyPromoCredit,
		ApplyPromoCreditAmount: req.ApplyPromoCreditAmount,
	})
	utility.AssertError(err, "Server Error")
	if err == nil && _interface.Context().Get(ctx).IsAdminPortalCall {
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
	var pendingCryptoSub *detail.SubscriptionDetail
	if createRes.Subscription != nil && createRes.Subscription.UserId > 0 {
		one := query.GetLatestCreateOrProcessingSubscriptionByUserId(ctx, createRes.Subscription.UserId, _interface.GetMerchantId(ctx), createRes.Plan.ProductId)
		if one != nil {
			gateway := query.GetGatewayById(ctx, one.GatewayId)
			if gateway.GatewayType == consts.GatewayTypeCrypto {
				pendingCryptoSub, _ = detail2.SubscriptionDetail(ctx, one.SubscriptionId)
			}
		}
	}
	token, err = jwt.CreatePortalToken(jwt.TOKENTYPEUSER, createRes.User.MerchantId, createRes.User.Id, createRes.User.Email, createRes.User.Language)
	// longer checkout token then 86400=1day
	utility.Assert(jwt.PutAuthTokenToCacheWithExpire(ctx, token, fmt.Sprintf("User#%d", createRes.User.Id), 86400), "Cache Error")
	g.RequestFromCtx(ctx).Cookie.Set("__UniBee.user.token", token)
	jwt.AppendRequestCookieWithToken(ctx, token)
	return &subscription.CreateRes{
		OtherPendingCryptoSubscription: pendingCryptoSub,
		Subscription:                   createRes.Subscription,
		User:                           createRes.User,
		Paid:                           createRes.Paid,
		Link:                           createRes.Link,
		Token:                          token,
	}, nil
}

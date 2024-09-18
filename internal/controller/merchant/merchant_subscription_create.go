package merchant

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"unibee/api/bean/detail"
	"unibee/internal/consts"
	_interface "unibee/internal/interface"
	"unibee/internal/logic/jwt"
	"unibee/internal/logic/operation_log"
	"unibee/internal/logic/subscription/service"
	user2 "unibee/internal/logic/user"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/merchant/subscription"
)

func (c *ControllerSubscription) Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error) {
	token := ""
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
		token, err = jwt.CreatePortalToken(jwt.TOKENTYPEUSER, user.MerchantId, user.Id, user.Email, user.Language)
		utility.AssertError(err, "Server Error")
		utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", user.Id)), "Cache Error")
		g.RequestFromCtx(ctx).Cookie.Set("__UniBee.user.token", token)
		jwt.AppendRequestCookieWithToken(ctx, token)
	} else if req.UserId == 0 && len(req.Email) > 0 {
		user, err := user2.QueryOrCreateUser(ctx, &user2.NewUserInternalReq{
			ExternalUserId: req.ExternalUserId,
			Email:          req.Email,
			MerchantId:     _interface.GetMerchantId(ctx),
		})
		utility.AssertError(err, "Server Error")
		req.UserId = user.Id
		token, err = jwt.CreatePortalToken(jwt.TOKENTYPEUSER, user.MerchantId, user.Id, user.Email, user.Language)
		utility.AssertError(err, "Server Error")
		utility.Assert(jwt.PutAuthTokenToCache(ctx, token, fmt.Sprintf("User#%d", user.Id)), "Cache Error")
		g.RequestFromCtx(ctx).Cookie.Set("__UniBee.user.token", token)
		jwt.AppendRequestCookieWithToken(ctx, token)
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
	var pendingCryptoSub *detail.SubscriptionDetail
	if createRes.Subscription != nil && createRes.Subscription.UserId > 0 {
		one := query.GetLatestCreateOrProcessingSubscriptionByUserId(ctx, createRes.Subscription.UserId, _interface.GetMerchantId(ctx), createRes.Plan.ProductId)
		if one != nil {
			gateway := query.GetGatewayById(ctx, one.GatewayId)
			if gateway.GatewayType == consts.GatewayTypeCrypto {
				pendingCryptoSub, _ = service.SubscriptionDetail(ctx, one.SubscriptionId)
			}
		}
	}
	return &subscription.CreateRes{
		OtherPendingCryptoSubscription: pendingCryptoSub,
		Subscription:                   createRes.Subscription,
		User:                           createRes.User,
		Paid:                           createRes.Paid,
		Link:                           createRes.Link,
		Token:                          token,
	}, nil
}

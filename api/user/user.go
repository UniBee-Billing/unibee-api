// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package user

import (
	"context"
	
	"go-oversea-pay/api/user/auth"
	"go-oversea-pay/api/user/invoice"
	"go-oversea-pay/api/user/payment"
	"go-oversea-pay/api/user/plan"
	"go-oversea-pay/api/user/profile"
	"go-oversea-pay/api/user/subscription"
	"go-oversea-pay/api/user/vat"
)

type IUserAuth interface {
	Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error)
	LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error)
	LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error)
	Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error)
	RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error)
}

type IUserInvoice interface {
	SubscriptionInvoiceList(ctx context.Context, req *invoice.SubscriptionInvoiceListReq) (res *invoice.SubscriptionInvoiceListRes, err error)
}

type IUserPayment interface {
	TimeLineList(ctx context.Context, req *payment.TimeLineListReq) (res *payment.TimeLineListRes, err error)
}

type IUserPlan interface {
	SubscriptionPlanList(ctx context.Context, req *plan.SubscriptionPlanListReq) (res *plan.SubscriptionPlanListRes, err error)
}

type IUserProfile interface {
	Profile(ctx context.Context, req *profile.ProfileReq) (res *profile.ProfileRes, err error)
	Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error)
	ProfileUpdate(ctx context.Context, req *profile.ProfileUpdateReq) (res *profile.ProfileUpdateRes, err error)
}

type IUserSubscription interface {
	SubscriptionDetail(ctx context.Context, req *subscription.SubscriptionDetailReq) (res *subscription.SubscriptionDetailRes, err error)
	SubscriptionPayCheck(ctx context.Context, req *subscription.SubscriptionPayCheckReq) (res *subscription.SubscriptionPayCheckRes, err error)
	SubscriptionChannels(ctx context.Context, req *subscription.SubscriptionChannelsReq) (res *subscription.SubscriptionChannelsRes, err error)
	SubscriptionCreatePreview(ctx context.Context, req *subscription.SubscriptionCreatePreviewReq) (res *subscription.SubscriptionCreatePreviewRes, err error)
	SubscriptionCreate(ctx context.Context, req *subscription.SubscriptionCreateReq) (res *subscription.SubscriptionCreateRes, err error)
	SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq) (res *subscription.SubscriptionUpdatePreviewRes, err error)
	SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq) (res *subscription.SubscriptionUpdateRes, err error)
	SubscriptionList(ctx context.Context, req *subscription.SubscriptionListReq) (res *subscription.SubscriptionListRes, err error)
	SubscriptionCancel(ctx context.Context, req *subscription.SubscriptionCancelReq) (res *subscription.SubscriptionCancelRes, err error)
	SubscriptionUpdateCancelAtPeriodEnd(ctx context.Context, req *subscription.SubscriptionUpdateCancelAtPeriodEndReq) (res *subscription.SubscriptionUpdateCancelAtPeriodEndRes, err error)
	SubscriptionUpdateCancelLastCancelAtPeriodEnd(ctx context.Context, req *subscription.SubscriptionUpdateCancelLastCancelAtPeriodEndReq) (res *subscription.SubscriptionUpdateCancelLastCancelAtPeriodEndRes, err error)
	SubscriptionSuspend(ctx context.Context, req *subscription.SubscriptionSuspendReq) (res *subscription.SubscriptionSuspendRes, err error)
	SubscriptionResume(ctx context.Context, req *subscription.SubscriptionResumeReq) (res *subscription.SubscriptionResumeRes, err error)
	SubscriptionTimeLineList(ctx context.Context, req *subscription.SubscriptionTimeLineListReq) (res *subscription.SubscriptionTimeLineListRes, err error)
}

type IUserVat interface {
	CountryVatList(ctx context.Context, req *vat.CountryVatListReq) (res *vat.CountryVatListRes, err error)
	NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error)
}



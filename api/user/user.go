// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package user

import (
	"context"
	
	"unibee/api/user/auth"
	"unibee/api/user/gateway"
	"unibee/api/user/invoice"
	"unibee/api/user/merchant"
	"unibee/api/user/payment"
	"unibee/api/user/plan"
	"unibee/api/user/profile"
	"unibee/api/user/subscription"
	"unibee/api/user/vat"
)

type IUserAuth interface {
	Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error)
	SessionLogin(ctx context.Context, req *auth.SessionLoginReq) (res *auth.SessionLoginRes, err error)
	LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error)
	LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error)
	PasswordForgetOtp(ctx context.Context, req *auth.PasswordForgetOtpReq) (res *auth.PasswordForgetOtpRes, err error)
	PasswordForgetOtpVerify(ctx context.Context, req *auth.PasswordForgetOtpVerifyReq) (res *auth.PasswordForgetOtpVerifyRes, err error)
	Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error)
	RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error)
}

type IUserGateway interface {
	List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error)
}

type IUserInvoice interface {
	List(ctx context.Context, req *invoice.ListReq) (res *invoice.ListRes, err error)
	Detail(ctx context.Context, req *invoice.DetailReq) (res *invoice.DetailRes, err error)
}

type IUserMerchant interface {
	Get(ctx context.Context, req *merchant.GetReq) (res *merchant.GetRes, err error)
}

type IUserPayment interface {
	ItemList(ctx context.Context, req *payment.ItemListReq) (res *payment.ItemListRes, err error)
	MethodList(ctx context.Context, req *payment.MethodListReq) (res *payment.MethodListRes, err error)
	MethodGet(ctx context.Context, req *payment.MethodGetReq) (res *payment.MethodGetRes, err error)
	MethodNew(ctx context.Context, req *payment.MethodNewReq) (res *payment.MethodNewRes, err error)
	MethodDelete(ctx context.Context, req *payment.MethodDeleteReq) (res *payment.MethodDeleteRes, err error)
	New(ctx context.Context, req *payment.NewReq) (res *payment.NewRes, err error)
	Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error)
	TimeLineList(ctx context.Context, req *payment.TimeLineListReq) (res *payment.TimeLineListRes, err error)
}

type IUserPlan interface {
	List(ctx context.Context, req *plan.ListReq) (res *plan.ListRes, err error)
}

type IUserProfile interface {
	Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error)
	Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error)
	Update(ctx context.Context, req *profile.UpdateReq) (res *profile.UpdateRes, err error)
	PasswordReset(ctx context.Context, req *profile.PasswordResetReq) (res *profile.PasswordResetRes, err error)
	ChangeGateway(ctx context.Context, req *profile.ChangeGatewayReq) (res *profile.ChangeGatewayRes, err error)
}

type IUserSubscription interface {
	UserCurrentSubscriptionDetail(ctx context.Context, req *subscription.UserCurrentSubscriptionDetailReq) (res *subscription.UserCurrentSubscriptionDetailRes, err error)
	Detail(ctx context.Context, req *subscription.DetailReq) (res *subscription.DetailRes, err error)
	PayCheck(ctx context.Context, req *subscription.PayCheckReq) (res *subscription.PayCheckRes, err error)
	CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error)
	Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error)
	UpdatePreview(ctx context.Context, req *subscription.UpdatePreviewReq) (res *subscription.UpdatePreviewRes, err error)
	Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error)
	List(ctx context.Context, req *subscription.ListReq) (res *subscription.ListRes, err error)
	Cancel(ctx context.Context, req *subscription.CancelReq) (res *subscription.CancelRes, err error)
	CancelAtPeriodEnd(ctx context.Context, req *subscription.CancelAtPeriodEndReq) (res *subscription.CancelAtPeriodEndRes, err error)
	CancelLastCancelAtPeriodEnd(ctx context.Context, req *subscription.CancelLastCancelAtPeriodEndReq) (res *subscription.CancelLastCancelAtPeriodEndRes, err error)
	Suspend(ctx context.Context, req *subscription.SuspendReq) (res *subscription.SuspendRes, err error)
	Resume(ctx context.Context, req *subscription.ResumeReq) (res *subscription.ResumeRes, err error)
	ChangeGateway(ctx context.Context, req *subscription.ChangeGatewayReq) (res *subscription.ChangeGatewayRes, err error)
	TimeLineList(ctx context.Context, req *subscription.TimeLineListReq) (res *subscription.TimeLineListRes, err error)
	OnetimeAddonNew(ctx context.Context, req *subscription.OnetimeAddonNewReq) (res *subscription.OnetimeAddonNewRes, err error)
	OnetimeAddonList(ctx context.Context, req *subscription.OnetimeAddonListReq) (res *subscription.OnetimeAddonListRes, err error)
	MarkWireTransferPaid(ctx context.Context, req *subscription.MarkWireTransferPaidReq) (res *subscription.MarkWireTransferPaidRes, err error)
	UserPendingCryptoSubscriptionDetail(ctx context.Context, req *subscription.UserPendingCryptoSubscriptionDetailReq) (res *subscription.UserPendingCryptoSubscriptionDetailRes, err error)
}

type IUserVat interface {
	CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error)
	NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error)
}



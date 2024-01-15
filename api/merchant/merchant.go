// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package merchant

import (
	"context"
	
	"go-oversea-pay/api/merchant/auth"
	"go-oversea-pay/api/merchant/invoice"
	"go-oversea-pay/api/merchant/oss"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/api/merchant/profile"
	"go-oversea-pay/api/merchant/subscription"
	"go-oversea-pay/api/merchant/vat"
	"go-oversea-pay/api/merchant/webhook"
)

type IMerchantAuth interface {
	Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error)
	LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error)
	LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error)
	Logout(ctx context.Context, req *auth.LogoutReq) (res *auth.LogoutRes, err error)
	Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error)
	RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error)
}

type IMerchantInvoice interface {
	SubscriptionInvoicePdfGenerate(ctx context.Context, req *invoice.SubscriptionInvoicePdfGenerateReq) (res *invoice.SubscriptionInvoicePdfGenerateRes, err error)
	SubscriptionInvoiceSendEmail(ctx context.Context, req *invoice.SubscriptionInvoiceSendEmailReq) (res *invoice.SubscriptionInvoiceSendEmailRes, err error)
	SubscriptionInvoiceList(ctx context.Context, req *invoice.SubscriptionInvoiceListReq) (res *invoice.SubscriptionInvoiceListRes, err error)
	NewInvoiceCreate(ctx context.Context, req *invoice.NewInvoiceCreateReq) (res *invoice.NewInvoiceCreateRes, err error)
}

type IMerchantOss interface {
	FileUpload(ctx context.Context, req *oss.FileUploadReq) (res *oss.FileUploadRes, err error)
}

type IMerchantPlan interface {
	SubscriptionPlanCreate(ctx context.Context, req *plan.SubscriptionPlanCreateReq) (res *plan.SubscriptionPlanCreateRes, err error)
	SubscriptionPlanEdit(ctx context.Context, req *plan.SubscriptionPlanEditReq) (res *plan.SubscriptionPlanEditRes, err error)
	SubscriptionPlanAddonsBinding(ctx context.Context, req *plan.SubscriptionPlanAddonsBindingReq) (res *plan.SubscriptionPlanAddonsBindingRes, err error)
	SubscriptionPlanList(ctx context.Context, req *plan.SubscriptionPlanListReq) (res *plan.SubscriptionPlanListRes, err error)
	SubscriptionPlanChannelTransferAndActivate(ctx context.Context, req *plan.SubscriptionPlanChannelTransferAndActivateReq) (res *plan.SubscriptionPlanChannelTransferAndActivateRes, err error)
	SubscriptionPlanChannelActivate(ctx context.Context, req *plan.SubscriptionPlanChannelActivateReq) (res *plan.SubscriptionPlanChannelActivateRes, err error)
	SubscriptionPlanChannelDeactivate(ctx context.Context, req *plan.SubscriptionPlanChannelDeactivateReq) (res *plan.SubscriptionPlanChannelDeactivateRes, err error)
	SubscriptionPlanDetail(ctx context.Context, req *plan.SubscriptionPlanDetailReq) (res *plan.SubscriptionPlanDetailRes, err error)
	SubscriptionPlanExpire(ctx context.Context, req *plan.SubscriptionPlanExpireReq) (res *plan.SubscriptionPlanExpireRes, err error)
}

type IMerchantProfile interface {
	Profile(ctx context.Context, req *profile.ProfileReq) (res *profile.ProfileRes, err error)
}

type IMerchantSubscription interface {
	SubscriptionDetail(ctx context.Context, req *subscription.SubscriptionDetailReq) (res *subscription.SubscriptionDetailRes, err error)
	SubscriptionList(ctx context.Context, req *subscription.SubscriptionListReq) (res *subscription.SubscriptionListRes, err error)
	SubscriptionCancel(ctx context.Context, req *subscription.SubscriptionCancelReq) (res *subscription.SubscriptionCancelRes, err error)
	SubscriptionUpdateCancelAtPeriodEnd(ctx context.Context, req *subscription.SubscriptionUpdateCancelAtPeriodEndReq) (res *subscription.SubscriptionUpdateCancelAtPeriodEndRes, err error)
	SubscriptionUpdateCancelLastCancelAtPeriodEnd(ctx context.Context, req *subscription.SubscriptionUpdateCancelLastCancelAtPeriodEndReq) (res *subscription.SubscriptionUpdateCancelLastCancelAtPeriodEndRes, err error)
	SubscriptionSuspend(ctx context.Context, req *subscription.SubscriptionSuspendReq) (res *subscription.SubscriptionSuspendRes, err error)
	SubscriptionResume(ctx context.Context, req *subscription.SubscriptionResumeReq) (res *subscription.SubscriptionResumeRes, err error)
	SubscriptionAddNewTrialStart(ctx context.Context, req *subscription.SubscriptionAddNewTrialStartReq) (res *subscription.SubscriptionAddNewTrialStartRes, err error)
	SubscriptionUpdatePreview(ctx context.Context, req *subscription.SubscriptionUpdatePreviewReq) (res *subscription.SubscriptionUpdatePreviewRes, err error)
	SubscriptionUpdate(ctx context.Context, req *subscription.SubscriptionUpdateReq) (res *subscription.SubscriptionUpdateRes, err error)
}

type IMerchantVat interface {
	SetupVatGateway(ctx context.Context, req *vat.SetupVatGatewayReq) (res *vat.SetupVatGatewayRes, err error)
	CountryVatList(ctx context.Context, req *vat.CountryVatListReq) (res *vat.CountryVatListRes, err error)
}

type IMerchantWebhook interface {
	SubscriptionWebhookCheckAndSetup(ctx context.Context, req *webhook.SubscriptionWebhookCheckAndSetupReq) (res *webhook.SubscriptionWebhookCheckAndSetupRes, err error)
}



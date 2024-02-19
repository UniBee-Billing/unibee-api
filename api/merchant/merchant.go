// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package merchant

import (
	"context"
	
	"unibee-api/api/merchant/auth"
	"unibee-api/api/merchant/balance"
	"unibee-api/api/merchant/email"
	"unibee-api/api/merchant/gateway"
	"unibee-api/api/merchant/invoice"
	"unibee-api/api/merchant/merchantinfo"
	"unibee-api/api/merchant/oss"
	"unibee-api/api/merchant/payment"
	"unibee-api/api/merchant/plan"
	"unibee-api/api/merchant/profile"
	"unibee-api/api/merchant/search"
	"unibee-api/api/merchant/subscription"
	"unibee-api/api/merchant/user"
	"unibee-api/api/merchant/vat"
	"unibee-api/api/merchant/webhook"
)

type IMerchantAuth interface {
	Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error)
	LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error)
	LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error)
	PasswordForgetOtp(ctx context.Context, req *auth.PasswordForgetOtpReq) (res *auth.PasswordForgetOtpRes, err error)
	PasswordForgetOtpVerify(ctx context.Context, req *auth.PasswordForgetOtpVerifyReq) (res *auth.PasswordForgetOtpVerifyRes, err error)
	Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error)
	RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error)
}

type IMerchantBalance interface {
	DetailQuery(ctx context.Context, req *balance.DetailQueryReq) (res *balance.DetailQueryRes, err error)
	UserDetailQuery(ctx context.Context, req *balance.UserDetailQueryReq) (res *balance.UserDetailQueryRes, err error)
}

type IMerchantEmail interface {
	MerchantEmailTemplateList(ctx context.Context, req *email.MerchantEmailTemplateListReq) (res *email.MerchantEmailTemplateListRes, err error)
	MerchantEmailTemplateUpdate(ctx context.Context, req *email.MerchantEmailTemplateUpdateReq) (res *email.MerchantEmailTemplateUpdateRes, err error)
	MerchantEmailTemplateSetDefault(ctx context.Context, req *email.MerchantEmailTemplateSetDefaultReq) (res *email.MerchantEmailTemplateSetDefaultRes, err error)
	MerchantEmailTemplateActivate(ctx context.Context, req *email.MerchantEmailTemplateActivateReq) (res *email.MerchantEmailTemplateActivateRes, err error)
	MerchantEmailTemplateDeactivate(ctx context.Context, req *email.MerchantEmailTemplateDeactivateReq) (res *email.MerchantEmailTemplateDeactivateRes, err error)
}

type IMerchantGateway interface {
	CheckAndSetup(ctx context.Context, req *gateway.CheckAndSetupReq) (res *gateway.CheckAndSetupRes, err error)
}

type IMerchantInvoice interface {
	SubscriptionInvoicePdfGenerate(ctx context.Context, req *invoice.SubscriptionInvoicePdfGenerateReq) (res *invoice.SubscriptionInvoicePdfGenerateRes, err error)
	SubscriptionInvoiceSendEmail(ctx context.Context, req *invoice.SubscriptionInvoiceSendEmailReq) (res *invoice.SubscriptionInvoiceSendEmailRes, err error)
	SubscriptionInvoiceDetail(ctx context.Context, req *invoice.SubscriptionInvoiceDetailReq) (res *invoice.SubscriptionInvoiceDetailRes, err error)
	SubscriptionInvoiceList(ctx context.Context, req *invoice.SubscriptionInvoiceListReq) (res *invoice.SubscriptionInvoiceListRes, err error)
	NewInvoiceCreate(ctx context.Context, req *invoice.NewInvoiceCreateReq) (res *invoice.NewInvoiceCreateRes, err error)
	NewInvoiceEdit(ctx context.Context, req *invoice.NewInvoiceEditReq) (res *invoice.NewInvoiceEditRes, err error)
	DeletePendingInvoice(ctx context.Context, req *invoice.DeletePendingInvoiceReq) (res *invoice.DeletePendingInvoiceRes, err error)
	ProcessInvoiceForPay(ctx context.Context, req *invoice.ProcessInvoiceForPayReq) (res *invoice.ProcessInvoiceForPayRes, err error)
	CancelProcessingInvoice(ctx context.Context, req *invoice.CancelProcessingInvoiceReq) (res *invoice.CancelProcessingInvoiceRes, err error)
	NewInvoiceRefund(ctx context.Context, req *invoice.NewInvoiceRefundReq) (res *invoice.NewInvoiceRefundRes, err error)
}

type IMerchantMerchantinfo interface {
	MerchantInfo(ctx context.Context, req *merchantinfo.MerchantInfoReq) (res *merchantinfo.MerchantInfoRes, err error)
	MerchantInfoUpdate(ctx context.Context, req *merchantinfo.MerchantInfoUpdateReq) (res *merchantinfo.MerchantInfoUpdateRes, err error)
}

type IMerchantOss interface {
	FileUpload(ctx context.Context, req *oss.FileUploadReq) (res *oss.FileUploadRes, err error)
}

type IMerchantPayment interface {
	TimeLineList(ctx context.Context, req *payment.TimeLineListReq) (res *payment.TimeLineListRes, err error)
}

type IMerchantPlan interface {
	SubscriptionPlanCreate(ctx context.Context, req *plan.SubscriptionPlanCreateReq) (res *plan.SubscriptionPlanCreateRes, err error)
	SubscriptionPlanEdit(ctx context.Context, req *plan.SubscriptionPlanEditReq) (res *plan.SubscriptionPlanEditRes, err error)
	SubscriptionPlanAddonsBinding(ctx context.Context, req *plan.SubscriptionPlanAddonsBindingReq) (res *plan.SubscriptionPlanAddonsBindingRes, err error)
	SubscriptionPlanList(ctx context.Context, req *plan.SubscriptionPlanListReq) (res *plan.SubscriptionPlanListRes, err error)
	SubscriptionPlanChannelTransferAndActivate(ctx context.Context, req *plan.SubscriptionPlanChannelTransferAndActivateReq) (res *plan.SubscriptionPlanChannelTransferAndActivateRes, err error)
	SubscriptionPlanChannelActivate(ctx context.Context, req *plan.SubscriptionPlanChannelActivateReq) (res *plan.SubscriptionPlanChannelActivateRes, err error)
	SubscriptionPlanChannelDeactivate(ctx context.Context, req *plan.SubscriptionPlanChannelDeactivateReq) (res *plan.SubscriptionPlanChannelDeactivateRes, err error)
	SubscriptionPlanPublish(ctx context.Context, req *plan.SubscriptionPlanPublishReq) (res *plan.SubscriptionPlanPublishRes, err error)
	SubscriptionPlanUnPublish(ctx context.Context, req *plan.SubscriptionPlanUnPublishReq) (res *plan.SubscriptionPlanUnPublishRes, err error)
	SubscriptionPlanDetail(ctx context.Context, req *plan.SubscriptionPlanDetailReq) (res *plan.SubscriptionPlanDetailRes, err error)
	SubscriptionPlanExpire(ctx context.Context, req *plan.SubscriptionPlanExpireReq) (res *plan.SubscriptionPlanExpireRes, err error)
}

type IMerchantProfile interface {
	Profile(ctx context.Context, req *profile.ProfileReq) (res *profile.ProfileRes, err error)
	Logout(ctx context.Context, req *profile.LogoutReq) (res *profile.LogoutRes, err error)
	PasswordReset(ctx context.Context, req *profile.PasswordResetReq) (res *profile.PasswordResetRes, err error)
}

type IMerchantSearch interface {
	Search(ctx context.Context, req *search.SearchReq) (res *search.SearchRes, err error)
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
	UserSubscriptionDetail(ctx context.Context, req *subscription.UserSubscriptionDetailReq) (res *subscription.UserSubscriptionDetailRes, err error)
	SubscriptionTimeLineList(ctx context.Context, req *subscription.SubscriptionTimeLineListReq) (res *subscription.SubscriptionTimeLineListRes, err error)
	SubscriptionMerchantPendingUpdateList(ctx context.Context, req *subscription.SubscriptionMerchantPendingUpdateListReq) (res *subscription.SubscriptionMerchantPendingUpdateListRes, err error)
	SubscriptionNewAdminNote(ctx context.Context, req *subscription.SubscriptionNewAdminNoteReq) (res *subscription.SubscriptionNewAdminNoteRes, err error)
	SubscriptionAdminNoteList(ctx context.Context, req *subscription.SubscriptionAdminNoteListReq) (res *subscription.SubscriptionAdminNoteListRes, err error)
}

type IMerchantUser interface {
	List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error)
	Get(ctx context.Context, req *user.GetReq) (res *user.GetRes, err error)
	Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error)
	UserProfileUpdate(ctx context.Context, req *user.UserProfileUpdateReq) (res *user.UserProfileUpdateRes, err error)
}

type IMerchantVat interface {
	SetupVatGateway(ctx context.Context, req *vat.SetupVatGatewayReq) (res *vat.SetupVatGatewayRes, err error)
	CountryVatList(ctx context.Context, req *vat.CountryVatListReq) (res *vat.CountryVatListRes, err error)
}

type IMerchantWebhook interface {
	EventList(ctx context.Context, req *webhook.EventListReq) (res *webhook.EventListRes, err error)
	EndpointList(ctx context.Context, req *webhook.EndpointListReq) (res *webhook.EndpointListRes, err error)
	EndpointLogList(ctx context.Context, req *webhook.EndpointLogListReq) (res *webhook.EndpointLogListRes, err error)
	NewEndpoint(ctx context.Context, req *webhook.NewEndpointReq) (res *webhook.NewEndpointRes, err error)
	UpdateEndpoint(ctx context.Context, req *webhook.UpdateEndpointReq) (res *webhook.UpdateEndpointRes, err error)
	DeleteEndpoint(ctx context.Context, req *webhook.DeleteEndpointReq) (res *webhook.DeleteEndpointRes, err error)
}



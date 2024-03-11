// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package merchant

import (
	"context"
	
	"unibee/api/merchant/auth"
	"unibee/api/merchant/balance"
	"unibee/api/merchant/email"
	"unibee/api/merchant/gateway"
	"unibee/api/merchant/invoice"
	"unibee/api/merchant/member"
	"unibee/api/merchant/metric"
	"unibee/api/merchant/oss"
	"unibee/api/merchant/payment"
	"unibee/api/merchant/plan"
	"unibee/api/merchant/profile"
	"unibee/api/merchant/search"
	"unibee/api/merchant/session"
	"unibee/api/merchant/subscription"
	"unibee/api/merchant/user"
	"unibee/api/merchant/vat"
	"unibee/api/merchant/webhook"
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
	GatewaySetup(ctx context.Context, req *email.GatewaySetupReq) (res *email.GatewaySetupRes, err error)
	TemplateList(ctx context.Context, req *email.TemplateListReq) (res *email.TemplateListRes, err error)
	TemplateUpdate(ctx context.Context, req *email.TemplateUpdateReq) (res *email.TemplateUpdateRes, err error)
	TemplateSetDefault(ctx context.Context, req *email.TemplateSetDefaultReq) (res *email.TemplateSetDefaultRes, err error)
	TemplateActivate(ctx context.Context, req *email.TemplateActivateReq) (res *email.TemplateActivateRes, err error)
	TemplateDeactivate(ctx context.Context, req *email.TemplateDeactivateReq) (res *email.TemplateDeactivateRes, err error)
}

type IMerchantGateway interface {
	List(ctx context.Context, req *gateway.ListReq) (res *gateway.ListRes, err error)
	Setup(ctx context.Context, req *gateway.SetupReq) (res *gateway.SetupRes, err error)
	Edit(ctx context.Context, req *gateway.EditReq) (res *gateway.EditRes, err error)
	SetupWebhook(ctx context.Context, req *gateway.SetupWebhookReq) (res *gateway.SetupWebhookRes, err error)
}

type IMerchantInvoice interface {
	PdfGenerate(ctx context.Context, req *invoice.PdfGenerateReq) (res *invoice.PdfGenerateRes, err error)
	SendEmail(ctx context.Context, req *invoice.SendEmailReq) (res *invoice.SendEmailRes, err error)
	Detail(ctx context.Context, req *invoice.DetailReq) (res *invoice.DetailRes, err error)
	List(ctx context.Context, req *invoice.ListReq) (res *invoice.ListRes, err error)
	New(ctx context.Context, req *invoice.NewReq) (res *invoice.NewRes, err error)
	Edit(ctx context.Context, req *invoice.EditReq) (res *invoice.EditRes, err error)
	Delete(ctx context.Context, req *invoice.DeleteReq) (res *invoice.DeleteRes, err error)
	Finish(ctx context.Context, req *invoice.FinishReq) (res *invoice.FinishRes, err error)
	Cancel(ctx context.Context, req *invoice.CancelReq) (res *invoice.CancelRes, err error)
	Refund(ctx context.Context, req *invoice.RefundReq) (res *invoice.RefundRes, err error)
}

type IMerchantMember interface {
	Profile(ctx context.Context, req *member.ProfileReq) (res *member.ProfileRes, err error)
	Logout(ctx context.Context, req *member.LogoutReq) (res *member.LogoutRes, err error)
	PasswordReset(ctx context.Context, req *member.PasswordResetReq) (res *member.PasswordResetRes, err error)
}

type IMerchantMetric interface {
	List(ctx context.Context, req *metric.ListReq) (res *metric.ListRes, err error)
	New(ctx context.Context, req *metric.NewReq) (res *metric.NewRes, err error)
	Edit(ctx context.Context, req *metric.EditReq) (res *metric.EditRes, err error)
	Delete(ctx context.Context, req *metric.DeleteReq) (res *metric.DeleteRes, err error)
	Detail(ctx context.Context, req *metric.DetailReq) (res *metric.DetailRes, err error)
	NewEvent(ctx context.Context, req *metric.NewEventReq) (res *metric.NewEventRes, err error)
	DeleteEvent(ctx context.Context, req *metric.DeleteEventReq) (res *metric.DeleteEventRes, err error)
	NewPlanLimit(ctx context.Context, req *metric.NewPlanLimitReq) (res *metric.NewPlanLimitRes, err error)
	EditPlanLimit(ctx context.Context, req *metric.EditPlanLimitReq) (res *metric.EditPlanLimitRes, err error)
	DeletePlanLimit(ctx context.Context, req *metric.DeletePlanLimitReq) (res *metric.DeletePlanLimitRes, err error)
	UserMetric(ctx context.Context, req *metric.UserMetricReq) (res *metric.UserMetricRes, err error)
}

type IMerchantOss interface {
	FileUpload(ctx context.Context, req *oss.FileUploadReq) (res *oss.FileUploadRes, err error)
}

type IMerchantPayment interface {
	Cancel(ctx context.Context, req *payment.CancelReq) (res *payment.CancelRes, err error)
	RefundCancel(ctx context.Context, req *payment.RefundCancelReq) (res *payment.RefundCancelRes, err error)
	Capture(ctx context.Context, req *payment.CaptureReq) (res *payment.CaptureRes, err error)
	MethodList(ctx context.Context, req *payment.MethodListReq) (res *payment.MethodListRes, err error)
	New(ctx context.Context, req *payment.NewReq) (res *payment.NewRes, err error)
	Detail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error)
	List(ctx context.Context, req *payment.ListReq) (res *payment.ListRes, err error)
	NewPaymentRefund(ctx context.Context, req *payment.NewPaymentRefundReq) (res *payment.NewPaymentRefundRes, err error)
	RefundDetail(ctx context.Context, req *payment.RefundDetailReq) (res *payment.RefundDetailRes, err error)
	RefundList(ctx context.Context, req *payment.RefundListReq) (res *payment.RefundListRes, err error)
	TimeLineList(ctx context.Context, req *payment.TimeLineListReq) (res *payment.TimeLineListRes, err error)
}

type IMerchantPlan interface {
	New(ctx context.Context, req *plan.NewReq) (res *plan.NewRes, err error)
	Edit(ctx context.Context, req *plan.EditReq) (res *plan.EditRes, err error)
	AddonsBinding(ctx context.Context, req *plan.AddonsBindingReq) (res *plan.AddonsBindingRes, err error)
	List(ctx context.Context, req *plan.ListReq) (res *plan.ListRes, err error)
	Activate(ctx context.Context, req *plan.ActivateReq) (res *plan.ActivateRes, err error)
	Publish(ctx context.Context, req *plan.PublishReq) (res *plan.PublishRes, err error)
	UnPublish(ctx context.Context, req *plan.UnPublishReq) (res *plan.UnPublishRes, err error)
	Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error)
	Expire(ctx context.Context, req *plan.ExpireReq) (res *plan.ExpireRes, err error)
	Delete(ctx context.Context, req *plan.DeleteReq) (res *plan.DeleteRes, err error)
}

type IMerchantProfile interface {
	Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error)
	Update(ctx context.Context, req *profile.UpdateReq) (res *profile.UpdateRes, err error)
}

type IMerchantSearch interface {
	Search(ctx context.Context, req *search.SearchReq) (res *search.SearchRes, err error)
}

type IMerchantSession interface {
	New(ctx context.Context, req *session.NewReq) (res *session.NewRes, err error)
}

type IMerchantSubscription interface {
	Detail(ctx context.Context, req *subscription.DetailReq) (res *subscription.DetailRes, err error)
	List(ctx context.Context, req *subscription.ListReq) (res *subscription.ListRes, err error)
	Cancel(ctx context.Context, req *subscription.CancelReq) (res *subscription.CancelRes, err error)
	CancelAtPeriodEnd(ctx context.Context, req *subscription.CancelAtPeriodEndReq) (res *subscription.CancelAtPeriodEndRes, err error)
	CancelLastCancelAtPeriodEnd(ctx context.Context, req *subscription.CancelLastCancelAtPeriodEndReq) (res *subscription.CancelLastCancelAtPeriodEndRes, err error)
	Suspend(ctx context.Context, req *subscription.SuspendReq) (res *subscription.SuspendRes, err error)
	Resume(ctx context.Context, req *subscription.ResumeReq) (res *subscription.ResumeRes, err error)
	AddNewTrialStart(ctx context.Context, req *subscription.AddNewTrialStartReq) (res *subscription.AddNewTrialStartRes, err error)
	UpdatePreview(ctx context.Context, req *subscription.UpdatePreviewReq) (res *subscription.UpdatePreviewRes, err error)
	Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error)
	UserSubscriptionDetail(ctx context.Context, req *subscription.UserSubscriptionDetailReq) (res *subscription.UserSubscriptionDetailRes, err error)
	TimeLineList(ctx context.Context, req *subscription.TimeLineListReq) (res *subscription.TimeLineListRes, err error)
	PendingUpdateList(ctx context.Context, req *subscription.PendingUpdateListReq) (res *subscription.PendingUpdateListRes, err error)
	NewAdminNote(ctx context.Context, req *subscription.NewAdminNoteReq) (res *subscription.NewAdminNoteRes, err error)
	AdminNoteList(ctx context.Context, req *subscription.AdminNoteListReq) (res *subscription.AdminNoteListRes, err error)
}

type IMerchantUser interface {
	List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error)
	Get(ctx context.Context, req *user.GetReq) (res *user.GetRes, err error)
	Frozen(ctx context.Context, req *user.FrozenReq) (res *user.FrozenRes, err error)
	Release(ctx context.Context, req *user.ReleaseReq) (res *user.ReleaseRes, err error)
	Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error)
	Update(ctx context.Context, req *user.UpdateReq) (res *user.UpdateRes, err error)
}

type IMerchantVat interface {
	SetupGateway(ctx context.Context, req *vat.SetupGatewayReq) (res *vat.SetupGatewayRes, err error)
	CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error)
}

type IMerchantWebhook interface {
	EventList(ctx context.Context, req *webhook.EventListReq) (res *webhook.EventListRes, err error)
	EndpointList(ctx context.Context, req *webhook.EndpointListReq) (res *webhook.EndpointListRes, err error)
	EndpointLogList(ctx context.Context, req *webhook.EndpointLogListReq) (res *webhook.EndpointLogListRes, err error)
	NewEndpoint(ctx context.Context, req *webhook.NewEndpointReq) (res *webhook.NewEndpointRes, err error)
	UpdateEndpoint(ctx context.Context, req *webhook.UpdateEndpointReq) (res *webhook.UpdateEndpointRes, err error)
	DeleteEndpoint(ctx context.Context, req *webhook.DeleteEndpointReq) (res *webhook.DeleteEndpointRes, err error)
}



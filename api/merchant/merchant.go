// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package merchant

import (
	"context"
	
	"unibee/api/merchant/auth"
	"unibee/api/merchant/balance"
	"unibee/api/merchant/discount"
	"unibee/api/merchant/email"
	"unibee/api/merchant/gateway"
	"unibee/api/merchant/invoice"
	"unibee/api/merchant/member"
	"unibee/api/merchant/metric"
	"unibee/api/merchant/oss"
	"unibee/api/merchant/payment"
	"unibee/api/merchant/plan"
	"unibee/api/merchant/product"
	"unibee/api/merchant/profile"
	"unibee/api/merchant/role"
	"unibee/api/merchant/search"
	"unibee/api/merchant/session"
	"unibee/api/merchant/subscription"
	"unibee/api/merchant/task"
	"unibee/api/merchant/track"
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

type IMerchantDiscount interface {
	List(ctx context.Context, req *discount.ListReq) (res *discount.ListRes, err error)
	Detail(ctx context.Context, req *discount.DetailReq) (res *discount.DetailRes, err error)
	New(ctx context.Context, req *discount.NewReq) (res *discount.NewRes, err error)
	Edit(ctx context.Context, req *discount.EditReq) (res *discount.EditRes, err error)
	Delete(ctx context.Context, req *discount.DeleteReq) (res *discount.DeleteRes, err error)
	Activate(ctx context.Context, req *discount.ActivateReq) (res *discount.ActivateRes, err error)
	Deactivate(ctx context.Context, req *discount.DeactivateReq) (res *discount.DeactivateRes, err error)
	UserDiscountList(ctx context.Context, req *discount.UserDiscountListReq) (res *discount.UserDiscountListRes, err error)
	PlanApplyPreview(ctx context.Context, req *discount.PlanApplyPreviewReq) (res *discount.PlanApplyPreviewRes, err error)
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
	EditCountryConfig(ctx context.Context, req *gateway.EditCountryConfigReq) (res *gateway.EditCountryConfigRes, err error)
	SetupWebhook(ctx context.Context, req *gateway.SetupWebhookReq) (res *gateway.SetupWebhookRes, err error)
	WireTransferSetup(ctx context.Context, req *gateway.WireTransferSetupReq) (res *gateway.WireTransferSetupRes, err error)
	WireTransferEdit(ctx context.Context, req *gateway.WireTransferEditReq) (res *gateway.WireTransferEditRes, err error)
}

type IMerchantInvoice interface {
	PdfGenerate(ctx context.Context, req *invoice.PdfGenerateReq) (res *invoice.PdfGenerateRes, err error)
	PdfUpdate(ctx context.Context, req *invoice.PdfUpdateReq) (res *invoice.PdfUpdateRes, err error)
	SendEmail(ctx context.Context, req *invoice.SendEmailReq) (res *invoice.SendEmailRes, err error)
	ReconvertCryptoAndSend(ctx context.Context, req *invoice.ReconvertCryptoAndSendReq) (res *invoice.ReconvertCryptoAndSendRes, err error)
	Detail(ctx context.Context, req *invoice.DetailReq) (res *invoice.DetailRes, err error)
	List(ctx context.Context, req *invoice.ListReq) (res *invoice.ListRes, err error)
	New(ctx context.Context, req *invoice.NewReq) (res *invoice.NewRes, err error)
	Edit(ctx context.Context, req *invoice.EditReq) (res *invoice.EditRes, err error)
	Delete(ctx context.Context, req *invoice.DeleteReq) (res *invoice.DeleteRes, err error)
	Finish(ctx context.Context, req *invoice.FinishReq) (res *invoice.FinishRes, err error)
	Cancel(ctx context.Context, req *invoice.CancelReq) (res *invoice.CancelRes, err error)
	Refund(ctx context.Context, req *invoice.RefundReq) (res *invoice.RefundRes, err error)
	MarkRefund(ctx context.Context, req *invoice.MarkRefundReq) (res *invoice.MarkRefundRes, err error)
	MarkWireTransferSuccess(ctx context.Context, req *invoice.MarkWireTransferSuccessReq) (res *invoice.MarkWireTransferSuccessRes, err error)
	MarkRefundInvoiceSuccess(ctx context.Context, req *invoice.MarkRefundInvoiceSuccessReq) (res *invoice.MarkRefundInvoiceSuccessRes, err error)
}

type IMerchantMember interface {
	Profile(ctx context.Context, req *member.ProfileReq) (res *member.ProfileRes, err error)
	Logout(ctx context.Context, req *member.LogoutReq) (res *member.LogoutRes, err error)
	PasswordReset(ctx context.Context, req *member.PasswordResetReq) (res *member.PasswordResetRes, err error)
	List(ctx context.Context, req *member.ListReq) (res *member.ListRes, err error)
	UpdateMemberRole(ctx context.Context, req *member.UpdateMemberRoleReq) (res *member.UpdateMemberRoleRes, err error)
	NewMember(ctx context.Context, req *member.NewMemberReq) (res *member.NewMemberRes, err error)
	Frozen(ctx context.Context, req *member.FrozenReq) (res *member.FrozenRes, err error)
	Release(ctx context.Context, req *member.ReleaseReq) (res *member.ReleaseRes, err error)
	OperationLogList(ctx context.Context, req *member.OperationLogListReq) (res *member.OperationLogListRes, err error)
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
	ItemList(ctx context.Context, req *payment.ItemListReq) (res *payment.ItemListRes, err error)
	MethodList(ctx context.Context, req *payment.MethodListReq) (res *payment.MethodListRes, err error)
	MethodGet(ctx context.Context, req *payment.MethodGetReq) (res *payment.MethodGetRes, err error)
	MethodNew(ctx context.Context, req *payment.MethodNewReq) (res *payment.MethodNewRes, err error)
	MethodDelete(ctx context.Context, req *payment.MethodDeleteReq) (res *payment.MethodDeleteRes, err error)
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
	Copy(ctx context.Context, req *plan.CopyReq) (res *plan.CopyRes, err error)
	Activate(ctx context.Context, req *plan.ActivateReq) (res *plan.ActivateRes, err error)
	Publish(ctx context.Context, req *plan.PublishReq) (res *plan.PublishRes, err error)
	UnPublish(ctx context.Context, req *plan.UnPublishReq) (res *plan.UnPublishRes, err error)
	Detail(ctx context.Context, req *plan.DetailReq) (res *plan.DetailRes, err error)
	Expire(ctx context.Context, req *plan.ExpireReq) (res *plan.ExpireRes, err error)
	Delete(ctx context.Context, req *plan.DeleteReq) (res *plan.DeleteRes, err error)
}

type IMerchantProduct interface {
	New(ctx context.Context, req *product.NewReq) (res *product.NewRes, err error)
	Edit(ctx context.Context, req *product.EditReq) (res *product.EditRes, err error)
	List(ctx context.Context, req *product.ListReq) (res *product.ListRes, err error)
	Copy(ctx context.Context, req *product.CopyReq) (res *product.CopyRes, err error)
	Activate(ctx context.Context, req *product.ActivateReq) (res *product.ActivateRes, err error)
	Inactive(ctx context.Context, req *product.InactiveReq) (res *product.InactiveRes, err error)
	Detail(ctx context.Context, req *product.DetailReq) (res *product.DetailRes, err error)
	Delete(ctx context.Context, req *product.DeleteReq) (res *product.DeleteRes, err error)
}

type IMerchantProfile interface {
	Get(ctx context.Context, req *profile.GetReq) (res *profile.GetRes, err error)
	Update(ctx context.Context, req *profile.UpdateReq) (res *profile.UpdateRes, err error)
	CountryConfigList(ctx context.Context, req *profile.CountryConfigListReq) (res *profile.CountryConfigListRes, err error)
	EditCountryConfig(ctx context.Context, req *profile.EditCountryConfigReq) (res *profile.EditCountryConfigRes, err error)
	NewApiKey(ctx context.Context, req *profile.NewApiKeyReq) (res *profile.NewApiKeyRes, err error)
}

type IMerchantRole interface {
	List(ctx context.Context, req *role.ListReq) (res *role.ListRes, err error)
	New(ctx context.Context, req *role.NewReq) (res *role.NewRes, err error)
	Edit(ctx context.Context, req *role.EditReq) (res *role.EditRes, err error)
	Delete(ctx context.Context, req *role.DeleteReq) (res *role.DeleteRes, err error)
}

type IMerchantSearch interface {
	Search(ctx context.Context, req *search.SearchReq) (res *search.SearchRes, err error)
}

type IMerchantSession interface {
	New(ctx context.Context, req *session.NewReq) (res *session.NewRes, err error)
}

type IMerchantSubscription interface {
	Config(ctx context.Context, req *subscription.ConfigReq) (res *subscription.ConfigRes, err error)
	ConfigUpdate(ctx context.Context, req *subscription.ConfigUpdateReq) (res *subscription.ConfigUpdateRes, err error)
	Detail(ctx context.Context, req *subscription.DetailReq) (res *subscription.DetailRes, err error)
	UserPendingCryptoSubscriptionDetail(ctx context.Context, req *subscription.UserPendingCryptoSubscriptionDetailReq) (res *subscription.UserPendingCryptoSubscriptionDetailRes, err error)
	List(ctx context.Context, req *subscription.ListReq) (res *subscription.ListRes, err error)
	Cancel(ctx context.Context, req *subscription.CancelReq) (res *subscription.CancelRes, err error)
	CancelAtPeriodEnd(ctx context.Context, req *subscription.CancelAtPeriodEndReq) (res *subscription.CancelAtPeriodEndRes, err error)
	CancelLastCancelAtPeriodEnd(ctx context.Context, req *subscription.CancelLastCancelAtPeriodEndReq) (res *subscription.CancelLastCancelAtPeriodEndRes, err error)
	Suspend(ctx context.Context, req *subscription.SuspendReq) (res *subscription.SuspendRes, err error)
	Resume(ctx context.Context, req *subscription.ResumeReq) (res *subscription.ResumeRes, err error)
	ChangeGateway(ctx context.Context, req *subscription.ChangeGatewayReq) (res *subscription.ChangeGatewayRes, err error)
	AddNewTrialStart(ctx context.Context, req *subscription.AddNewTrialStartReq) (res *subscription.AddNewTrialStartRes, err error)
	Renew(ctx context.Context, req *subscription.RenewReq) (res *subscription.RenewRes, err error)
	CreatePreview(ctx context.Context, req *subscription.CreatePreviewReq) (res *subscription.CreatePreviewRes, err error)
	Create(ctx context.Context, req *subscription.CreateReq) (res *subscription.CreateRes, err error)
	UpdatePreview(ctx context.Context, req *subscription.UpdatePreviewReq) (res *subscription.UpdatePreviewRes, err error)
	Update(ctx context.Context, req *subscription.UpdateReq) (res *subscription.UpdateRes, err error)
	UserSubscriptionDetail(ctx context.Context, req *subscription.UserSubscriptionDetailReq) (res *subscription.UserSubscriptionDetailRes, err error)
	TimeLineList(ctx context.Context, req *subscription.TimeLineListReq) (res *subscription.TimeLineListRes, err error)
	PendingUpdateList(ctx context.Context, req *subscription.PendingUpdateListReq) (res *subscription.PendingUpdateListRes, err error)
	NewAdminNote(ctx context.Context, req *subscription.NewAdminNoteReq) (res *subscription.NewAdminNoteRes, err error)
	ActiveTemporarily(ctx context.Context, req *subscription.ActiveTemporarilyReq) (res *subscription.ActiveTemporarilyRes, err error)
	AdminNoteList(ctx context.Context, req *subscription.AdminNoteListReq) (res *subscription.AdminNoteListRes, err error)
	OnetimeAddonNew(ctx context.Context, req *subscription.OnetimeAddonNewReq) (res *subscription.OnetimeAddonNewRes, err error)
	OnetimeAddonList(ctx context.Context, req *subscription.OnetimeAddonListReq) (res *subscription.OnetimeAddonListRes, err error)
	NewPayment(ctx context.Context, req *subscription.NewPaymentReq) (res *subscription.NewPaymentRes, err error)
}

type IMerchantTask interface {
	List(ctx context.Context, req *task.ListReq) (res *task.ListRes, err error)
	ExportColumnList(ctx context.Context, req *task.ExportColumnListReq) (res *task.ExportColumnListRes, err error)
	New(ctx context.Context, req *task.NewReq) (res *task.NewRes, err error)
	NewImport(ctx context.Context, req *task.NewImportReq) (res *task.NewImportRes, err error)
	NewTemplate(ctx context.Context, req *task.NewTemplateReq) (res *task.NewTemplateRes, err error)
	EditTemplate(ctx context.Context, req *task.EditTemplateReq) (res *task.EditTemplateRes, err error)
	DeleteTemplate(ctx context.Context, req *task.DeleteTemplateReq) (res *task.DeleteTemplateRes, err error)
	ExportTemplateList(ctx context.Context, req *task.ExportTemplateListReq) (res *task.ExportTemplateListRes, err error)
}

type IMerchantTrack interface {
	SetupSegment(ctx context.Context, req *track.SetupSegmentReq) (res *track.SetupSegmentRes, err error)
}

type IMerchantUser interface {
	New(ctx context.Context, req *user.NewReq) (res *user.NewRes, err error)
	List(ctx context.Context, req *user.ListReq) (res *user.ListRes, err error)
	Get(ctx context.Context, req *user.GetReq) (res *user.GetRes, err error)
	Frozen(ctx context.Context, req *user.FrozenReq) (res *user.FrozenRes, err error)
	Release(ctx context.Context, req *user.ReleaseReq) (res *user.ReleaseRes, err error)
	Search(ctx context.Context, req *user.SearchReq) (res *user.SearchRes, err error)
	Update(ctx context.Context, req *user.UpdateReq) (res *user.UpdateRes, err error)
	ChangeGateway(ctx context.Context, req *user.ChangeGatewayReq) (res *user.ChangeGatewayRes, err error)
}

type IMerchantVat interface {
	SetupGateway(ctx context.Context, req *vat.SetupGatewayReq) (res *vat.SetupGatewayRes, err error)
	InitDefaultGateway(ctx context.Context, req *vat.InitDefaultGatewayReq) (res *vat.InitDefaultGatewayRes, err error)
	CountryList(ctx context.Context, req *vat.CountryListReq) (res *vat.CountryListRes, err error)
	NumberValidate(ctx context.Context, req *vat.NumberValidateReq) (res *vat.NumberValidateRes, err error)
}

type IMerchantWebhook interface {
	EventList(ctx context.Context, req *webhook.EventListReq) (res *webhook.EventListRes, err error)
	EndpointList(ctx context.Context, req *webhook.EndpointListReq) (res *webhook.EndpointListRes, err error)
	EndpointLogList(ctx context.Context, req *webhook.EndpointLogListReq) (res *webhook.EndpointLogListRes, err error)
	ResendWebhook(ctx context.Context, req *webhook.ResendWebhookReq) (res *webhook.ResendWebhookRes, err error)
	NewEndpoint(ctx context.Context, req *webhook.NewEndpointReq) (res *webhook.NewEndpointRes, err error)
	UpdateEndpoint(ctx context.Context, req *webhook.UpdateEndpointReq) (res *webhook.UpdateEndpointRes, err error)
	DeleteEndpoint(ctx context.Context, req *webhook.DeleteEndpointReq) (res *webhook.DeleteEndpointRes, err error)
}



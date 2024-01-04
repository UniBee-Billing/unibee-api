// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package merchant

import (
	"context"

	"go-oversea-pay/api/merchant/auth"
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/api/merchant/profile"
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

type IMerchantWebhook interface {
	SubscriptionWebhookCheckAndSetup(ctx context.Context, req *webhook.SubscriptionWebhookCheckAndSetupReq) (res *webhook.SubscriptionWebhookCheckAndSetupRes, err error)
}

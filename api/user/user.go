// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package user

import (
	"context"
	
	"go-oversea-pay/api/user/auth"
	"go-oversea-pay/api/user/plan"
	"go-oversea-pay/api/user/profile"
	"go-oversea-pay/api/user/subscription"
)

type IUserAuth interface {
	Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error)
	LoginOtp(ctx context.Context, req *auth.LoginOtpReq) (res *auth.LoginOtpRes, err error)
	LoginOtpVerify(ctx context.Context, req *auth.LoginOtpVerifyReq) (res *auth.LoginOtpVerifyRes, err error)
	Logout(ctx context.Context, req *auth.LogoutReq) (res *auth.LogoutRes, err error)
	Register(ctx context.Context, req *auth.RegisterReq) (res *auth.RegisterRes, err error)
	RegisterVerify(ctx context.Context, req *auth.RegisterVerifyReq) (res *auth.RegisterVerifyRes, err error)
}

type IUserPlan interface {
	SubscriptionPlanList(ctx context.Context, req *plan.SubscriptionPlanListReq) (res *plan.SubscriptionPlanListRes, err error)
}

type IUserProfile interface {
	Profile(ctx context.Context, req *profile.ProfileReq) (res *profile.ProfileRes, err error)
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
	SubscriptionSuspend(ctx context.Context, req *subscription.SubscriptionSuspendReq) (res *subscription.SubscriptionSuspendRes, err error)
	SubscriptionResume(ctx context.Context, req *subscription.SubscriptionResumeReq) (res *subscription.SubscriptionResumeRes, err error)
}



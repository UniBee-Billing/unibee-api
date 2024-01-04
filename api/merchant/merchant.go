// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package merchant

import (
	"context"
	
	"go-oversea-pay/api/merchant/plan"
	"go-oversea-pay/api/merchant/subscription"
	"go-oversea-pay/api/merchant/webhook"
)

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

type IMerchantSubscription interface {
	SubscriptionList(ctx context.Context, req *subscription.SubscriptionListReq) (res *subscription.SubscriptionListRes, err error)
}

type IMerchantWebhook interface {
	SubscriptionWebhookCheckAndSetup(ctx context.Context, req *webhook.SubscriptionWebhookCheckAndSetupReq) (res *webhook.SubscriptionWebhookCheckAndSetupRes, err error)
}



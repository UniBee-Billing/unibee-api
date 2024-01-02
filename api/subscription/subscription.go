// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package subscription

import (
	"context"
	
	"go-oversea-pay/api/subscription/v1"
)

type ISubscriptionV1 interface {
	SubscriptionChannels(ctx context.Context, req *v1.SubscriptionChannelsReq) (res *v1.SubscriptionChannelsRes, err error)
	SubscriptionPlanCreate(ctx context.Context, req *v1.SubscriptionPlanCreateReq) (res *v1.SubscriptionPlanCreateRes, err error)
	SubscriptionPlanEdit(ctx context.Context, req *v1.SubscriptionPlanEditReq) (res *v1.SubscriptionPlanEditRes, err error)
	SubscriptionPlanAddonsBinding(ctx context.Context, req *v1.SubscriptionPlanAddonsBindingReq) (res *v1.SubscriptionPlanAddonsBindingRes, err error)
	SubscriptionPlanList(ctx context.Context, req *v1.SubscriptionPlanListReq) (res *v1.SubscriptionPlanListRes, err error)
	SubscriptionPlanChannelTransferAndActivate(ctx context.Context, req *v1.SubscriptionPlanChannelTransferAndActivateReq) (res *v1.SubscriptionPlanChannelTransferAndActivateRes, err error)
	SubscriptionPlanChannelActivate(ctx context.Context, req *v1.SubscriptionPlanChannelActivateReq) (res *v1.SubscriptionPlanChannelActivateRes, err error)
	SubscriptionPlanChannelDeactivate(ctx context.Context, req *v1.SubscriptionPlanChannelDeactivateReq) (res *v1.SubscriptionPlanChannelDeactivateRes, err error)
	SubscriptionPlanDetail(ctx context.Context, req *v1.SubscriptionPlanDetailReq) (res *v1.SubscriptionPlanDetailRes, err error)
	SubscriptionPlanExpire(ctx context.Context, req *v1.SubscriptionPlanExpireReq) (res *v1.SubscriptionPlanExpireRes, err error)
	SubscriptionCreate(ctx context.Context, req *v1.SubscriptionCreateReq) (res *v1.SubscriptionCreateRes, err error)
	SubscriptionCancel(ctx context.Context, req *v1.SubscriptionCancelReq) (res *v1.SubscriptionCancelRes, err error)
	SubscriptionDetail(ctx context.Context, req *v1.SubscriptionDetailReq) (res *v1.SubscriptionDetailRes, err error)
	SubscriptionUpdate(ctx context.Context, req *v1.SubscriptionUpdateReq) (res *v1.SubscriptionUpdateRes, err error)
	SubscriptionWebhookCheckAndSetup(ctx context.Context, req *v1.SubscriptionWebhookCheckAndSetupReq) (res *v1.SubscriptionWebhookCheckAndSetupRes, err error)
}



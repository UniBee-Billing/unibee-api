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
	SubscriptionPlanChannelTransferAndActivate(ctx context.Context, req *v1.SubscriptionPlanChannelTransferAndActivateReq) (res *v1.SubscriptionPlanChannelTransferAndActivateRes, err error)
	SubscriptionPlanChannelActivate(ctx context.Context, req *v1.SubscriptionPlanChannelActivateReq) (res *v1.SubscriptionPlanChannelActivateRes, err error)
	SubscriptionPlanChannelDeactivate(ctx context.Context, req *v1.SubscriptionPlanChannelDeactivateReq) (res *v1.SubscriptionPlanChannelDeactivateRes, err error)
	SubscriptionPlanDetail(ctx context.Context, req *v1.SubscriptionPlanDetailReq) (res *v1.SubscriptionPlanDetailRes, err error)
	SubscriptionCreate(ctx context.Context, req *v1.SubscriptionCreateReq) (res *v1.SubscriptionCreateRes, err error)
}



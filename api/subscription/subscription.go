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
	SubscriptionPlanChannelTransfer(ctx context.Context, req *v1.SubscriptionPlanChannelTransferReq) (res *v1.SubscriptionPlanChannelTransferRes, err error)
	SubscriptionPlanChannelActive(ctx context.Context, req *v1.SubscriptionPlanChannelActiveReq) (res *v1.SubscriptionPlanChannelActiveRes, err error)
	SubscriptionPlanChannelInActive(ctx context.Context, req *v1.SubscriptionPlanChannelInActiveReq) (res *v1.SubscriptionPlanChannelInActiveRes, err error)
	SubscriptionPlanDetail(ctx context.Context, req *v1.SubscriptionPlanDetailReq) (res *v1.SubscriptionPlanDetailRes, err error)
	SubscriptionCreate(ctx context.Context, req *v1.SubscriptionCreateReq) (res *v1.SubscriptionCreateRes, err error)
}



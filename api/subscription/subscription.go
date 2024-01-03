// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package subscription

import (
	"context"
	
	"go-oversea-pay/api/subscription/v1"
)

type ISubscriptionV1 interface {
	SubscriptionDetail(ctx context.Context, req *v1.SubscriptionDetailReq) (res *v1.SubscriptionDetailRes, err error)
	SubscriptionChannels(ctx context.Context, req *v1.SubscriptionChannelsReq) (res *v1.SubscriptionChannelsRes, err error)
	SubscriptionCreate(ctx context.Context, req *v1.SubscriptionCreateReq) (res *v1.SubscriptionCreateRes, err error)
	SubscriptionCancel(ctx context.Context, req *v1.SubscriptionCancelReq) (res *v1.SubscriptionCancelRes, err error)
	SubscriptionUpdate(ctx context.Context, req *v1.SubscriptionUpdateReq) (res *v1.SubscriptionUpdateRes, err error)
}



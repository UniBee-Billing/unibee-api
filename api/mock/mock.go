// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package mock

import (
	"context"
	
	"go-oversea-pay/api/mock/v1"
)

type IMockV1 interface {
	Cancel(ctx context.Context, req *v1.CancelReq) (res *v1.CancelRes, err error)
	Capture(ctx context.Context, req *v1.CaptureReq) (res *v1.CaptureRes, err error)
	SamplePaymentNetherlands(ctx context.Context, req *v1.SamplePaymentNetherlandsReq) (res *v1.SamplePaymentNetherlandsRes, err error)
	DetailPay(ctx context.Context, req *v1.DetailPayReq) (res *v1.DetailPayRes, err error)
	Refund(ctx context.Context, req *v1.RefundReq) (res *v1.RefundRes, err error)
	SubscriptionChannels(ctx context.Context, req *v1.SubscriptionChannelsReq) (res *v1.SubscriptionChannelsRes, err error)
	SubscriptionPlanCreate(ctx context.Context, req *v1.SubscriptionPlanCreateReq) (res *v1.SubscriptionPlanCreateRes, err error)
	SubscriptionPlanChannelTransfer(ctx context.Context, req *v1.SubscriptionPlanChannelTransferReq) (res *v1.SubscriptionPlanChannelTransferRes, err error)
	SubscriptionPlanDetail(ctx context.Context, req *v1.SubscriptionPlanDetailReq) (res *v1.SubscriptionPlanDetailRes, err error)
}



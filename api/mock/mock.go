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
}



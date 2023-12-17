// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package mock

import (
	"context"
	
	"go-oversea-pay/api/mock/v1"
)

type IMockV1 interface {
	SamplePaymentNetherlands(ctx context.Context, req *v1.SamplePaymentNetherlandsReq) (res *v1.SamplePaymentNetherlandsRes, err error)
}



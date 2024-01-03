// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package out

import (
	"context"
	
	"go-oversea-pay/api/out/v1"
)

type IOutV1 interface {
	Cancels(ctx context.Context, req *v1.CancelsReq) (res *v1.CancelsRes, err error)
	Captures(ctx context.Context, req *v1.CapturesReq) (res *v1.CapturesRes, err error)
	Payments(ctx context.Context, req *v1.PaymentsReq) (res *v1.PaymentsRes, err error)
	PaymentMethods(ctx context.Context, req *v1.PaymentMethodsReq) (res *v1.PaymentMethodsRes, err error)
	PaymentDetails(ctx context.Context, req *v1.PaymentDetailsReq) (res *v1.PaymentDetailsRes, err error)
	DisableRecurringDetails(ctx context.Context, req *v1.DisableRecurringDetailsReq) (res *v1.DisableRecurringDetailsRes, err error)
	ListRecurringDetails(ctx context.Context, req *v1.ListRecurringDetailsReq) (res *v1.ListRecurringDetailsRes, err error)
	Refunds(ctx context.Context, req *v1.RefundsReq) (res *v1.RefundsRes, err error)
}



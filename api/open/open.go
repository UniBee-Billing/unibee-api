// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT. 
// =================================================================================

package open

import (
	"context"
	
	"go-oversea-pay/api/open/mock"
	"go-oversea-pay/api/open/payment"
)

type IOpenMock interface {
	Cancel(ctx context.Context, req *mock.CancelReq) (res *mock.CancelRes, err error)
	Capture(ctx context.Context, req *mock.CaptureReq) (res *mock.CaptureRes, err error)
	SamplePaymentNetherlands(ctx context.Context, req *mock.SamplePaymentNetherlandsReq) (res *mock.SamplePaymentNetherlandsRes, err error)
	DetailPay(ctx context.Context, req *mock.DetailPayReq) (res *mock.DetailPayRes, err error)
	Refund(ctx context.Context, req *mock.RefundReq) (res *mock.RefundRes, err error)
}

type IOpenPayment interface {
	Cancels(ctx context.Context, req *payment.CancelsReq) (res *payment.CancelsRes, err error)
	Captures(ctx context.Context, req *payment.CapturesReq) (res *payment.CapturesRes, err error)
	Payments(ctx context.Context, req *payment.PaymentsReq) (res *payment.PaymentsRes, err error)
	PaymentMethods(ctx context.Context, req *payment.PaymentMethodsReq) (res *payment.PaymentMethodsRes, err error)
	PaymentDetails(ctx context.Context, req *payment.PaymentDetailsReq) (res *payment.PaymentDetailsRes, err error)
	DisableRecurringDetails(ctx context.Context, req *payment.DisableRecurringDetailsReq) (res *payment.DisableRecurringDetailsRes, err error)
	ListRecurringDetails(ctx context.Context, req *payment.ListRecurringDetailsReq) (res *payment.ListRecurringDetailsRes, err error)
	Refunds(ctx context.Context, req *payment.RefundsReq) (res *payment.RefundsRes, err error)
}



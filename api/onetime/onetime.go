// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package onetime

import (
	"context"

	"unibee-api/api/onetime/mock"
	"unibee-api/api/onetime/payment"
)

type IOpenMock interface {
	Cancel(ctx context.Context, req *mock.CancelReq) (res *mock.CancelRes, err error)
	Capture(ctx context.Context, req *mock.CaptureReq) (res *mock.CaptureRes, err error)
	SamplePaymentNetherlands(ctx context.Context, req *mock.SamplePaymentNetherlandsReq) (res *mock.SamplePaymentNetherlandsRes, err error)
	DetailPay(ctx context.Context, req *mock.DetailPayReq) (res *mock.DetailPayRes, err error)
	MockMessageSend(ctx context.Context, req *mock.MockMessageSendReq) (res *mock.MockMessageSendRes, err error)
	Refund(ctx context.Context, req *mock.RefundReq) (res *mock.RefundRes, err error)
}

type IOpenPayment interface {
	Cancel(ctx context.Context, req *payment.CancelReq) (res *payment.CancelRes, err error)
	Capture(ctx context.Context, req *payment.CaptureReq) (res *payment.CaptureRes, err error)
	NewPayment(ctx context.Context, req *payment.NewPaymentReq) (res *payment.NewPaymentRes, err error)
	PaymentMethodList(ctx context.Context, req *payment.MethodListReq) (res *payment.MethodListRes, err error)
	PaymentDetail(ctx context.Context, req *payment.DetailReq) (res *payment.DetailRes, err error)
	DisableRecurringDetail(ctx context.Context, req *payment.DisableRecurringDetailReq) (res *payment.DisableRecurringDetailRes, err error)
	RecurringDetailList(ctx context.Context, req *payment.RecurringDetailListReq) (res *payment.RecurringDetailListRes, err error)
	NewPaymentRefund(ctx context.Context, req *payment.NewPaymentRefundReq) (res *payment.NewPaymentRefundRes, err error)
}

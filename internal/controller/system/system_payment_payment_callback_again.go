package system

import (
	"context"
	"unibee/internal/consts"
	"unibee/internal/logic/payment/callback"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/system/payment"
)

func (c *ControllerPayment) PaymentCallbackAgain(ctx context.Context, req *payment.PaymentCallbackAgainReq) (res *payment.PaymentCallbackAgainRes, err error) {
	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, one.InvoiceId)
	utility.Assert(invoice != nil, "invoice not found")
	if one.Status == consts.PaymentSuccess {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentSuccessCallback(ctx, one, invoice)
	} else if one.Status == consts.PaymentFailed {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentFailureCallback(ctx, one, invoice)
	} else if one.Status == consts.PaymentCancelled {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentCancelCallback(ctx, one, invoice)
	} else if one.Status == consts.PaymentCreated && one.AuthorizeStatus == consts.WaitingAuthorized {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentNeedAuthorisedCallback(ctx, one, invoice)
	}
	return &payment.PaymentCallbackAgainRes{}, nil
}

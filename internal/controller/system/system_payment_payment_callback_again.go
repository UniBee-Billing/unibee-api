package system

import (
	"context"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/internal/logic/payment/callback"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/utility"

	"go-oversea-pay/api/system/payment"
)

func (c *ControllerPayment) PaymentCallbackAgain(ctx context.Context, req *payment.PaymentCallbackAgainReq) (res *payment.PaymentCallbackAgainRes, err error) {
	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, one.InvoiceId)
	utility.Assert(invoice != nil, "invoice not found")
	if one.Status == consts.PAY_SUCCESS {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentSuccessCallback(ctx, one, invoice)
	} else if one.Status == consts.PAY_FAILED {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentFailureCallback(ctx, one, invoice)
	} else if one.Status == consts.PAY_CANCEL {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentCancelCallback(ctx, one, invoice)
	} else if one.Status == consts.TO_BE_PAID && one.AuthorizeStatus == consts.WAITING_AUTHORIZED {
		callback.GetPaymentCallbackServiceProvider(ctx, one.BizType).PaymentNeedAuthorisedCallback(ctx, one, invoice)
	}
	return &payment.PaymentCallbackAgainRes{}, nil
}
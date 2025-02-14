package system

import (
	"context"
	redismq "github.com/jackyang-hk/go-redismq"
	redismqcmd "unibee/internal/cmd/redismq"
	"unibee/internal/query"
	"unibee/utility"

	"unibee/api/system/payment"
)

func (c *ControllerPayment) PaymentGatewayCheck(ctx context.Context, req *payment.PaymentGatewayCheckReq) (res *payment.PaymentGatewayCheckRes, err error) {
	one := query.GetPaymentByPaymentId(ctx, req.PaymentId)
	utility.Assert(one != nil, "payment not found")
	invoice := query.GetInvoiceByInvoiceId(ctx, one.InvoiceId)
	utility.Assert(invoice != nil, "invoice not found")
	// send the payment status checker mq
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismqcmd.TopicPaymentChecker.Topic,
		Tag:        redismqcmd.TopicPaymentChecker.Tag,
		Body:       req.PaymentId,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	return &payment.PaymentGatewayCheckRes{}, nil
}

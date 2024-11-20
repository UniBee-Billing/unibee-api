package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	"time"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/invoice"
	"unibee/internal/logic/discount"
	"unibee/internal/query"
	"unibee/utility"
)

type InvoiceReversedListener struct {
}

func (t InvoiceReversedListener) GetTopic() string {
	return redismq2.TopicInvoiceReversed.Topic
}

func (t InvoiceReversedListener) GetTag() string {
	return redismq2.TopicInvoiceReversed.Tag
}

func (t InvoiceReversedListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "InvoiceReversedListener Receive Message:%s", utility.MarshalToJsonString(message))
	one := query.GetInvoiceByInvoiceId(ctx, message.Body)
	if one != nil {
		if len(one.DiscountCode) > 0 {
			discount.UpdateUserDiscountPaymentIdWhenInvoicePaid(ctx, one.InvoiceId, one.PaymentId)
		}
		one.Status = consts.InvoiceStatusReversed
		go func() {
			time.Sleep(300 * time.Millisecond)
			invoice.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_PAID, message.CustomData)
		}()
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewInvoiceReversedListener())
	fmt.Println("NewInvoiceReversedListener RegisterListener")
}

func NewInvoiceReversedListener() *InvoiceReversedListener {
	return &InvoiceReversedListener{}
}

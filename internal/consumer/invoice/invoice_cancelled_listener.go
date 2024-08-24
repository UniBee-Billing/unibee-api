package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/invoice"
	"unibee/internal/query"
	"unibee/utility"
)

type InvoiceCancelledListener struct {
}

func (t InvoiceCancelledListener) GetTopic() string {
	return redismq2.TopicInvoiceCancelled.Topic
}

func (t InvoiceCancelledListener) GetTag() string {
	return redismq2.TopicInvoiceCancelled.Tag
}

func (t InvoiceCancelledListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "InvoiceCancelledListener Receive Message:%s", utility.MarshalToJsonString(message))
	one := query.GetInvoiceByInvoiceId(ctx, message.Body)
	if one != nil {
		one.Status = consts.InvoiceStatusCancelled
		invoice.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_CANCELLED)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewInvoiceCancelledListener())
	fmt.Println("InvoiceCancelledListener RegisterListener")
}

func NewInvoiceCancelledListener() *InvoiceCancelledListener {
	return &InvoiceCancelledListener{}
}

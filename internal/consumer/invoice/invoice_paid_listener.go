package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/internal/consts"
	"unibee/internal/consumer/webhook/event"
	"unibee/internal/consumer/webhook/invoice"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type InvoicePaidListener struct {
}

func (t InvoicePaidListener) GetTopic() string {
	return redismq2.TopicInvoicePaid.Topic
}

func (t InvoicePaidListener) GetTag() string {
	return redismq2.TopicInvoicePaid.Tag
}

func (t InvoicePaidListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "InvoicePaidListener Receive Message:%s", utility.MarshalToJsonString(message))
	one := query.GetInvoiceByInvoiceId(ctx, message.Body)
	if one != nil {
		one.Status = consts.InvoiceStatusPaid
		invoice.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_PAID)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewInvoicePaidListener())
	fmt.Println("NewInvoicePaidListener RegisterListener")
}

func NewInvoicePaidListener() *InvoicePaidListener {
	return &InvoicePaidListener{}
}

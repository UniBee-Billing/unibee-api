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

type InvoiceCreateListener struct {
}

func (t InvoiceCreateListener) GetTopic() string {
	return redismq2.TopicInvoiceCreated.Topic
}

func (t InvoiceCreateListener) GetTag() string {
	return redismq2.TopicInvoiceCreated.Tag
}

func (t InvoiceCreateListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "InvoiceCreateListener Receive Message:%s", utility.MarshalToJsonString(message))
	one := query.GetInvoiceByInvoiceId(ctx, message.Body)
	if one != nil {
		one.Status = consts.InvoiceStatusPending
		invoice.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_CREATED)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewInvoiceCreateListener())
	fmt.Println("InvoiceCreateListener RegisterListener")
}

func NewInvoiceCreateListener() *InvoiceCreateListener {
	return &InvoiceCreateListener{}
}

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

type InvoiceFailedListener struct {
}

func (t InvoiceFailedListener) GetTopic() string {
	return redismq2.TopicInvoiceFailed.Topic
}

func (t InvoiceFailedListener) GetTag() string {
	return redismq2.TopicInvoiceFailed.Tag
}

func (t InvoiceFailedListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "InvoiceFailedListener Receive Message:%s", utility.MarshalToJsonString(message))
	one := query.GetInvoiceByInvoiceId(ctx, message.Body)
	if one != nil {
		one.Status = consts.InvoiceStatusFailed
		invoice.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_FAILED)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewInvoiceFailedListener())
	fmt.Println("InvoiceFailedListener RegisterListener")
}

func NewInvoiceFailedListener() *InvoiceFailedListener {
	return &InvoiceFailedListener{}
}

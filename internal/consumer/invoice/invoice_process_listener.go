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

type InvoiceProcessListener struct {
}

func (t InvoiceProcessListener) GetTopic() string {
	return redismq2.TopicInvoiceProcessed.Topic
}

func (t InvoiceProcessListener) GetTag() string {
	return redismq2.TopicInvoiceProcessed.Tag
}

func (t InvoiceProcessListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Debugf(ctx, "InvoiceProcessListener Receive Message:%s", utility.MarshalToJsonString(message))
	one := query.GetInvoiceByInvoiceId(ctx, message.Body)
	if one != nil {
		one.Status = consts.InvoiceStatusProcessing
		invoice.SendMerchantInvoiceWebhookBackground(one, event.UNIBEE_WEBHOOK_EVENT_INVOICE_PROCESS)
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewInvoiceProcessListener())
	fmt.Println("InvoiceProcessListener RegisterListener")
}

func NewInvoiceProcessListener() *InvoiceProcessListener {
	return &InvoiceProcessListener{}
}

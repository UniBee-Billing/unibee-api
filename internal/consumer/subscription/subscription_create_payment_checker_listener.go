package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq2 "unibee-api/internal/cmd/redismq"
	"unibee-api/internal/consts"
	"unibee-api/internal/logic/invoice/handler"
	"unibee-api/internal/logic/subscription/billingcycle/expire"
	"unibee-api/internal/query"
	"unibee-api/redismq"
	"unibee-api/utility"
)

type SubscriptionCreatePaymentCheckListener struct {
}

func (t SubscriptionCreatePaymentCheckListener) GetTopic() string {
	return redismq2.TopicSubscriptionCreatePaymentCheck.Topic
}

func (t SubscriptionCreatePaymentCheckListener) GetTag() string {
	return redismq2.TopicSubscriptionCreatePaymentCheck.Tag
}

func (t SubscriptionCreatePaymentCheckListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "SubscriptionCreatePaymentCheckListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)

	if gtime.Now().Timestamp()-sub.GmtCreate.Timestamp() >= 2*24*60*60 {
		//should expire sub
		err := expire.SubscriptionExpire(ctx, sub, "NotPayAfter48Hours")
		if err != nil {
			fmt.Printf("SubscriptionCreatePaymentCheckListener SubscriptionExpire Error:%s", err.Error())
		}
		return redismq.CommitMessage
	}

	// After 3min Not Pay Send Email
	if sub.Status == consts.SubStatusCreate && len(sub.LatestInvoiceId) > 0 {
		invoice := query.GetInvoiceByInvoiceId(ctx, sub.LatestInvoiceId)
		if invoice != nil && invoice.Status == consts.InvoiceStatusProcessing {
			err := handler.SendSubscriptionInvoiceEmailToUser(ctx, sub.LatestInvoiceId)
			_, _ = redismq.SendDelay(&redismq.Message{
				Topic: redismq2.TopicSubscriptionCreatePaymentCheck.Topic,
				Tag:   redismq2.TopicSubscriptionCreatePaymentCheck.Tag,
				Body:  sub.SubscriptionId,
			}, 24*60*60) //every day send util expire
			if err != nil {
				return redismq.CommitMessage
			}
		}
	}
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(NewSubscriptionCreatePaymentCheckListener())
	fmt.Println("SubscriptionCreatePaymentCheckListener RegisterListener")
}

func NewSubscriptionCreatePaymentCheckListener() *SubscriptionCreatePaymentCheckListener {
	return &SubscriptionCreatePaymentCheckListener{}
}

package subscription

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	redismq2 "go-oversea-pay/internal/cmd/redismq"
	"go-oversea-pay/internal/consts"
	sub2 "go-oversea-pay/internal/cronjob/sub"
	"go-oversea-pay/internal/logic/invoice/handler"
	"go-oversea-pay/internal/query"
	"go-oversea-pay/redismq"
	"go-oversea-pay/utility"
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
	fmt.Printf("SubscriptionCreatePaymentCheckListener Receive Message:%s", utility.MarshalToJsonString(message))
	sub := query.GetSubscriptionBySubscriptionId(ctx, message.Body)

	if gtime.Now().Timestamp()-sub.GmtCreate.Timestamp() > 50*60*60 {
		//should expire sub
		err := sub2.SubscriptionExpire(ctx, sub, "NotPayAfter48Hours")
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
			if err != nil {
				_, _ = redismq.SendDelay(&redismq.Message{
					Topic: redismq2.TopicSubscriptionCreatePaymentCheck.Topic,
					Tag:   redismq2.TopicSubscriptionCreatePaymentCheck.Tag,
					Body:  sub.SubscriptionId,
				}, 24*60*60)
				return redismq.CommitMessage
			} else {
				return redismq.ReconsumeLater
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

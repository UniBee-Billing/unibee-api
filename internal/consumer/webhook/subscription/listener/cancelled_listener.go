package listener

//
//import (
//	"context"
//	"github.com/gogf/gf/v2/frame/g"
//	redismq2 "unibee-api/internal/cmd/redismq"
//	"unibee-api/internal/consumer/webhook/event"
//	"unibee-api/internal/consumer/webhook/http"
//	message2 "unibee-api/internal/consumer/webhook/message"
//	"unibee-api/redismq"
//	"unibee-api/utility"
//)
//
//type CancelledListener struct {
//}
//
//func (t CancelledListener) GetTopic() string {
//	return redismq2.TopicMerchantWebhook.Topic
//}
//
//func (t CancelledListener) GetTag() string {
//	return event.MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CANCELLED
//}
//
//func (t CancelledListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
//	utility.Assert(len(message.Body) > 0, "body is nil")
//	utility.Assert(len(message.Body) != 0, "body length is 0")
//	g.Log().Infof(ctx, "Webhook_Subscription NewCancelledListener Receive Message:%s", utility.MarshalToJsonString(message))
//	var webhookMessage *message2.WebhookMessage
//	err := utility.UnmarshalFromJsonString(message.Body, &webhookMessage)
//
//	if err != nil {
//		g.Log().Infof(ctx, "Webhook_Subscription NewCancelledListener UnmarshalFromJsonString Error:%s", err.Error())
//		return redismq.ReconsumeLater
//	}
//
//	if http.SendWebhookRequest(ctx, webhookMessage.Url, webhookMessage.Data) {
//		return redismq.CommitMessage
//	}
//
//	return redismq.ReconsumeLater
//}
//
//func init() {
//	listener := NewCancelledListener()
//	redismq.RegisterListener(listener)
//	event.RegisterListenerEvent(listener)
//}
//
//func NewCancelledListener() *CancelledListener {
//	return &CancelledListener{}
//}

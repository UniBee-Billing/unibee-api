package message

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
	"strings"
	redismq2 "unibee/internal/cmd/redismq"
	event2 "unibee/internal/consumer/webhook/event"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"
)

type WebhookMessage struct {
	Id                uint64
	Event             event2.WebhookEvent
	EventId           string
	EndpointId        uint64
	Url               string
	MerchantId        uint64
	Data              *gjson.Json
	SequenceKey       string
	DependencyKey     string
	EndpointEventList string
}

func SendWebhookMessage(ctx context.Context, event event2.WebhookEvent, merchantId uint64, data *gjson.Json, sequenceKey string, dependencyKey string) {
	webhookMessage := &entity.MerchantWebhookMessage{
		MerchantId:      merchantId,
		WebhookEvent:    string(event),
		Data:            data.String(),
		WebsocketStatus: 10,
		CreateTime:      gtime.Now().Timestamp(),
	}
	insert, err := dao.MerchantWebhookMessage.Ctx(ctx).Data(webhookMessage).OmitNil().Insert(webhookMessage)
	utility.AssertError(err, "webhook message insert error")
	id, err := insert.LastInsertId()
	utility.AssertError(err, "webhook message LastInsertId error")
	webhookMessage.Id = uint64(id)

	eventId := utility.CreateEventId()

	{
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicInternalWebhook.Topic,
			Tag:   redismq2.TopicInternalWebhook.Tag,
			Body: utility.MarshalToJsonString(&WebhookMessage{
				Id:            webhookMessage.Id,
				Event:         event,
				EventId:       eventId,
				MerchantId:    merchantId,
				Data:          data,
				SequenceKey:   sequenceKey,
				DependencyKey: dependencyKey,
			}),
		})
	}

	utility.Assert(event2.WebhookEventInListeningEvents(event), fmt.Sprintf("Event:%s Not In Event List", event))
	list := query.GetMerchantWebhooksByMerchantId(ctx, merchantId)
	if list != nil {
		for _, merchantWebhook := range list {
			eventList := strings.Split(merchantWebhook.WebhookEvents, ",")
			if in(eventList, string(event)) {
				send, err := redismq.Send(&redismq.Message{
					Topic: redismq2.TopicMerchantWebhook.Topic,
					Tag:   redismq2.TopicMerchantWebhook.Tag,
					Body: utility.MarshalToJsonString(&WebhookMessage{
						Id:            webhookMessage.Id,
						Event:         event,
						EventId:       eventId,
						EndpointId:    merchantWebhook.Id,
						Url:           merchantWebhook.WebhookUrl,
						MerchantId:    merchantId,
						Data:          data,
						SequenceKey:   sequenceKey,
						DependencyKey: dependencyKey,
					}),
				})
				g.Log().Infof(ctx, "SendWebhookMessage event:%s, merchantWebhookUrl:%s send:%v err:%v", event, merchantWebhook.WebhookUrl, send, err)
			}
		}
	}
}

func in(strArray []string, target string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

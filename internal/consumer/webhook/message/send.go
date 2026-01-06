package message

import (
	"context"
	"fmt"
	"strings"
	redismq2 "unibee/internal/cmd/redismq"
	event2 "unibee/internal/consumer/webhook/event"
	dao "unibee/internal/dao/default"
	entity "unibee/internal/model/entity/default"
	"unibee/internal/query"
	"unibee/utility"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	redismq "github.com/jackyang-hk/go-redismq"
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
	MetaData          string
}

func SendWebhookMessage(ctx context.Context, event event2.WebhookEvent, merchantId uint64, data *gjson.Json, sequenceKey string, dependencyKey string, metadata map[string]interface{}) {
	var webhookMessageId uint64 = 0
	webhookMessage := &entity.MerchantWebhookMessage{
		MerchantId:      merchantId,
		WebhookEvent:    string(event),
		Data:            data.String(),
		WebsocketStatus: 10,
		CreateTime:      gtime.Now().Timestamp(),
	}
	insert, err := dao.MerchantWebhookMessage.Ctx(ctx).Data(webhookMessage).OmitNil().Insert(webhookMessage)
	if err != nil {
		g.Log().Errorf(ctx, "SendWebhookMessage insert merchant webhook message err:%s", err.Error())
	} else {
		id, err := insert.LastInsertId()
		if err != nil {
			g.Log().Errorf(ctx, "SendWebhookMessage insert merchant webhook message get LastInsertId err:%s", err.Error())
		} else {
			webhookMessage.Id = uint64(id)
			webhookMessageId = webhookMessage.Id
		}
	}

	if metadata != nil {
		if _, ok := metadata["Persistence"]; ok {
			if subId, ok := metadata["SubscriptionId"]; ok {
				persistenceOne := &entity.MerchantWebhookMessage{
					MerchantId:      merchantId,
					WebhookEvent:    string(event),
					Data:            data.String(),
					WebsocketStatus: 50,
					SubscriptionId:  fmt.Sprintf("%s", subId),
					CreateTime:      gtime.Now().Timestamp(),
				}
				_, err = dao.MerchantWebhookMessage.Ctx(ctx).Data(persistenceOne).OmitNil().Insert(persistenceOne)
				if err != nil {
					g.Log().Errorf(ctx, fmt.Sprintf("SendWebhookMessage, Persistence error:%s", err.Error()))
				}
			}
		}
	}

	eventId := utility.CreateEventId()

	{
		_, _ = redismq.Send(&redismq.Message{
			Topic: redismq2.TopicInternalWebhook.Topic,
			Tag:   redismq2.TopicInternalWebhook.Tag,
			Body: utility.MarshalToJsonString(&WebhookMessage{
				Id:            webhookMessageId,
				Event:         event,
				EventId:       eventId,
				MerchantId:    merchantId,
				Data:          data,
				SequenceKey:   sequenceKey,
				DependencyKey: dependencyKey,
				MetaData:      utility.MarshalToJsonString(metadata),
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
					Topic:                     redismq2.TopicMerchantWebhook.Topic,
					Tag:                       redismq2.TopicMerchantWebhook.Tag,
					ConsumerDelayMilliSeconds: 100,
					Body: utility.MarshalToJsonString(&WebhookMessage{
						Id:            webhookMessageId,
						Event:         event,
						EventId:       eventId,
						EndpointId:    merchantWebhook.Id,
						Url:           merchantWebhook.WebhookUrl,
						MerchantId:    merchantId,
						Data:          data,
						SequenceKey:   sequenceKey,
						DependencyKey: dependencyKey,
						MetaData:      utility.MarshalToJsonString(metadata),
					}),
				})
				if err != nil {
					g.Log().Errorf(ctx, "SendWebhookMessage event:%s, merchantWebhookUrl:%s send:%v err:%s", event, merchantWebhook.WebhookUrl, send, err.Error())
				} else {
					g.Log().Infof(ctx, "SendWebhookMessage event:%s, merchantWebhookUrl:%s send:%v", event, merchantWebhook.WebhookUrl, send)
				}
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

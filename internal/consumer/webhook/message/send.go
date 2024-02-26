package message

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	redismq2 "unibee/internal/cmd/redismq"
	event2 "unibee/internal/consumer/webhook/event"
	"unibee/internal/query"
	"unibee/redismq"
	"unibee/utility"
)

type WebhookMessage struct {
	Event      event2.MerchantWebhookEvent
	EndpointId uint64
	Url        string
	MerchantId uint64
	Data       *gjson.Json
}

func SendWebhookMessage(ctx context.Context, event event2.MerchantWebhookEvent, merchantId uint64, data *gjson.Json) {
	utility.Assert(event2.WebhookEventInListeningEvents(event), fmt.Sprintf("Event:%s Not In Event List", event))
	list := query.GetMerchantWebhooksByMerchantId(ctx, merchantId)
	if list != nil {
		for _, merchantWebhook := range list {
			if strings.Contains(merchantWebhook.WebhookEvents, string(event)) {
				send, err := redismq.Send(&redismq.Message{
					Topic: redismq2.TopicMerchantWebhook.Topic,
					Tag:   redismq2.TopicMerchantWebhook.Tag,
					Body: utility.MarshalToJsonString(&WebhookMessage{
						Event:      event,
						EndpointId: merchantWebhook.Id,
						Url:        merchantWebhook.WebhookUrl,
						MerchantId: merchantId,
						Data:       data,
					}),
				})
				g.Log().Infof(ctx, "SendWebhookMessage event:%s, merchantWebhookUrl:%s send:%v err:%v", event, merchantWebhook.WebhookUrl, send, err)
			}
		}
	}
}

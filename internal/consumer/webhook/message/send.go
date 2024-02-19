package message

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
	redismq2 "unibee-api/internal/cmd/redismq"
	event2 "unibee-api/internal/consumer/webhook/event"
	"unibee-api/internal/query"
	"unibee-api/redismq"
	"unibee-api/utility"
)

type WebhookMessage struct {
	Event      string
	EndpointId uint64
	Url        string
	MerchantId int64
	Data       *gjson.Json
}

func SendWebhookMessage(ctx context.Context, event string, merchantId int64, data *gjson.Json) {
	utility.Assert(event2.EventInListeningEvents(event), fmt.Sprintf("Event:%s Not In Event List", event))
	list := query.GetMerchantWebhooksByMerchantId(ctx, merchantId)
	if list != nil {
		for _, merchantWebhook := range list {
			if strings.Contains(merchantWebhook.WebhookEvents, event) {
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

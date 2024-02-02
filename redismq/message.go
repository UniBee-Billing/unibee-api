package redismq

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/redis/go-redis/v9"
	"go-oversea-pay/utility"
	"strings"
)

type Message struct {
	MessageId        string                 `json:"messageId" dc:"MessageId"`
	Topic            string                 `json:"topic" dc:"Topic"`
	Tag              string                 `json:"tag" dc:"Tag"`
	Body             []byte                 `json:"body" dc:"Body"`
	Key              string                 `json:"key" dc:"Key"`
	StartDeliverTime int64                  `json:"startDeliverTime" dc:"Send Time,0-No Delay，Mills"`
	ReconsumeTimes   int                    `json:"reconsumeTimes" dc:"Reconsume Count"`
	CustomData       map[string]interface{} `json:"customData" dc:"CustomData"`
	SendTime         int64                  `json:"sendTime" dc:"Sent Time"`
}

type MessageMetaData struct {
	StartDeliverTime int64                  `json:"startDeliverTime" dc:"Send Time,0-No Delay，Mills"`
	ReconsumeTimes   int                    `json:"reconsumeTimes" dc:"Reconsume Count"`
	CustomData       map[string]interface{} `json:"customData" dc:"CustomData"`
	Key              string                 `json:"key" dc:"Key"`
	SendTime         int64                  `json:"sendTime" dc:"SendTime"`
}

func NewRedisMQMessage(topicWrappper MQTopicEnum, body interface{}) *Message {
	return &Message{
		Topic:    topicWrappper.Topic,
		Tag:      topicWrappper.Tag,
		Body:     Serialize(body),
		SendTime: utility.CurrentTimeMillis(),
	}
}

func (message *Message) getUniqueKey() string {
	uniqueKey := message.CustomData["uniqueKey"].(string)
	if len(uniqueKey) == 0 && len(message.MessageId) > 0 {
		message.CustomData["uniqueKey"] = message.MessageId
		return message.MessageId
	} else {
		return uniqueKey
	}
}

func (message *Message) isBoardCastingMessage() bool {
	return strings.Compare(message.CustomData["messageModel"].(string), "BROADCASTING") == 0
}

func (message *Message) getDescription() string {
	return fmt.Sprintf("%s %s %s", message.MessageId, message.Topic, message.Tag)
}

func (message *Message) toStreamAddArgsValues(stream string) *redis.XAddArgs {
	metadata := MessageMetaData{
		StartDeliverTime: message.StartDeliverTime,
		ReconsumeTimes:   message.ReconsumeTimes,
		CustomData:       message.CustomData,
		Key:              message.Key,
		SendTime:         utility.CurrentTimeMillis(),
	}
	metajson, _ := gjson.Marshal(metadata)
	return &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"topic":    message.Topic,
			"tag":      message.Tag,
			"body":     message.Body,
			"metadata": string(metajson),
		},
	}
}

func (message *Message) paseStreamMessage(value map[string]interface{}) {
	message.Topic = value["topic"].(string)
	message.Tag = value["tag"].(string)
	message.Body = value["body"].([]byte)
	metadata := value["metadata"].(string)
	if len(metadata) > 0 {
		json, err := gjson.LoadJson(metadata, true)
		if err == nil {
			defer func() {
				if exception := recover(); exception != nil {
					fmt.Printf("Redismq paseStreamMessage panic error:%s\n", exception)
					return
				}
			}()
			message.ReconsumeTimes = json.Get("reconsumeTimes").Int()
			message.StartDeliverTime = json.Get("startDeliverTime").Int64()
			message.SendTime = json.Get("sendTime").Int64()
			message.CustomData = json.Get("sendTime").Map()
			message.Key = json.Get("key").String()
			message.getUniqueKey()
		}
	}
}

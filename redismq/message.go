package redismq

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/redis/go-redis/v9"
	"go-oversea-pay/utility"
	"strings"
)

type Message struct {
	MessageId        string                 `p:"messageId" dc:"消息Id"`
	Topic            string                 `p:"topic" dc:"消息Topic"`
	Tag              string                 `p:"tag" dc:"消息Tag"`
	Body             []byte                 `p:"body" dc:"消息Body"`
	Key              string                 `p:"key" dc:"消息Key"`
	StartDeliverTime int64                  `p:"startDeliverTime" dc:"消息发送时间,0-表示不延迟，单位毫秒"`
	ReconsumeTimes   int                    `p:"reconsumeTimes" dc:"消息消费时, 获取消息已经被重试消费的次数"`
	CustomData       map[string]interface{} `p:"customData" dc:"自定义数据"`
	SendTime         int64                  `p:"sendTime" dc:"消息发送时间"`
}

type MessageMetaData struct {
	StartDeliverTime int64                  `p:"startDeliverTime" dc:"消息发送时间,0-表示不延迟，单位毫秒"`
	ReconsumeTimes   int                    `p:"reconsumeTimes" dc:"消息消费时, 获取消息已经被重试消费的次数"`
	CustomData       map[string]interface{} `p:"customData" dc:"自定义数据"`
	Key              string                 `p:"key" dc:"消息Key"`
	SendTime         int64                  `p:"sendTime" dc:"消息发送时间"`
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
					fmt.Printf("redismq paseStreamMessage panic error:%s\n", exception)
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

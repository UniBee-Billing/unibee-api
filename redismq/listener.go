package redismq

import (
	"fmt"
	"strings"
)

type IMessageListener interface {
	GetTopic() string
	GetTag() string
	Consume(message *Message) Action
}

var listeners map[string]IMessageListener
var Topics []string

func Listeners() map[string]IMessageListener {
	if listeners == nil {
		listeners = make(map[string]IMessageListener)
	}
	return listeners
}

func isValidTopic(topic string) bool {
	return len(topic) > 0 && strings.Compare(topic, "*") != 0
}

func RegisterListener(i IMessageListener) {
	if i == nil {
		return
	}
	if Topics == nil {
		Topics = make([]string, 0, 100) //最多容纳 60 个 topic
	}
	if len(Topics) > 60 {
		fmt.Println("工程Topic过多，请聚合Topic")
		return
	}
	if !isValidTopic(i.GetTopic()) {
		fmt.Printf("redismq 注册默认消费者失败 无效topic:%s,已忽略\n", i.GetTopic())
		return
	}
	if Listeners()[GetMessageKey(i.GetTopic(), i.GetTag())] != nil {
		fmt.Printf("redismq 多个消费者%s,消费同一个消息:%s,已忽略\n", i, GetMessageKey(i.GetTopic(), i.GetTag()))
	} else {
		messageKey := GetMessageKey(i.GetTopic(), i.GetTag())
		Listeners()[messageKey] = i
		Topics = append(Topics, messageKey)
		fmt.Printf("redismq 注册MQ消费者 IMessageListener:%s,消费消息:%s\n", i, messageKey)
	}
}

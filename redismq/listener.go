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
		fmt.Println("Project Register Topic Too Much ，Merge Please")
		return
	}
	if !isValidTopic(i.GetTopic()) {
		fmt.Printf("Redismq Regist Default Consumer Invalid Topic:%s,Drop\n", i.GetTopic())
		return
	}
	if Listeners()[GetMessageKey(i.GetTopic(), i.GetTag())] != nil {
		fmt.Printf("Redismq Multi %s,Consumer On:%s,Drop\n", i, GetMessageKey(i.GetTopic(), i.GetTag()))
	} else {
		messageKey := GetMessageKey(i.GetTopic(), i.GetTag())
		Listeners()[messageKey] = i
		Topics = append(Topics, messageKey)
		fmt.Printf("Redismq Register IMessageListener:%s,Consumer:%s\n", i, messageKey)
	}
}

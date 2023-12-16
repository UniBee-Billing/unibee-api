package redismq

import (
	"fmt"
)

type IMessageChecker interface {
	GetTopic() string
	GetTag() string
	Checker(message *Message) TransactionStatus
}

var checkers map[string]IMessageChecker

func Checkers() map[string]IMessageChecker {
	if checkers == nil {
		checkers = make(map[string]IMessageChecker)
	}
	return checkers
}

func RegisterChecker(i IMessageChecker) {
	if i == nil {
		return
	}
	if Checkers()[GetMessageKey(i.GetTopic(), i.GetTag())] != nil {
		fmt.Printf("redismq 多个事务观察者%s,观察同一个消息:%s,已忽略\n", i, GetMessageKey(i.GetTopic(), i.GetTag()))
	} else {
		Checkers()[GetMessageKey(i.GetTopic(), i.GetTag())] = i
		fmt.Printf("redismq 注册MQ事务观察者 IMessageChecker:%s,观察消息:%s\n", i, GetMessageKey(i.GetTopic(), i.GetTag()))
	}
}

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
		fmt.Printf("Redismq Multi Register Transaction Consumer %s,Watch One Message:%s,Drop\n", i, GetMessageKey(i.GetTopic(), i.GetTag()))
	} else {
		Checkers()[GetMessageKey(i.GetTopic(), i.GetTag())] = i
		fmt.Printf("Redismq Regist Consumer IMessageChecker:%s,Watch:%s\n", i, GetMessageKey(i.GetTopic(), i.GetTag()))
	}
}

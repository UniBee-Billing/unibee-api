package redismq

type Action int

const (
	CommitMessage = iota
	ReconsumeLater
)

func (action Action) Description() string {
	switch action {
	case CommitMessage:
		return "CommitMessage"
	case ReconsumeLater:
		return "ReconsumeLater"
	default:
		return "CommitMessage"
	}
}

type TransactionStatus int

const (
	CommitTransaction = iota
	RollbackTransaction
	Unknown
)

func (status TransactionStatus) Description() string {
	switch status {
	case CommitTransaction:
		return "CommitTransaction"
	case RollbackTransaction:
		return "RollbackTransaction"
	case Unknown:
		return "Unknown"
	default:
		return "CommitTransaction"
	}
}

type MQTopicEnum struct {
	Topic       string
	Tag         string
	Description string
}

var (
	TopicBlank = MQTopicEnum{"blank", "blank", "redis blank test message"}
)

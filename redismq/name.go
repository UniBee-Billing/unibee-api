package redismq

import (
	"fmt"
	"strings"
)
import "github.com/google/uuid"

const (
	REDIS_MQ_QUEUE_KEY_PREFIX            = "MQ_QUEUE_LIST_"
	REDIS_MQ_BACKUP_QUEUE_KEY_PREFIX     = "MQ_BACKUP_QUEUE_LIST_"
	STREAM_NAME                          = "STREAM_"
	NAME_VERSION                         = "_V3"
	REDIS_MQ_TRANSACTION_PRE_QUEUE_KEY   = "MQ_TRANSACTION_PRE_QUEUE_LIST_"
	REDIS_MQ_DEATH_QUEUE_KEY             = "MQ_DEATH_QUEUE_LIST"
	REDIS_MQ_TRANSACTION_DEATH_QUEUE_KEY = "MQ_TRANSACTION_DEATH_QUEUE_LIST_"
)

func GetQueueName(topic string) string {
	return fmt.Sprintf("%s%s%s%s", REDIS_MQ_QUEUE_KEY_PREFIX, STREAM_NAME, topic, NAME_VERSION)
}

func getBackupQueueName(topic string) string {
	return fmt.Sprintf("%s%s%s%s", REDIS_MQ_BACKUP_QUEUE_KEY_PREFIX, STREAM_NAME, topic, NAME_VERSION)
}

func GenerateUniqueNo(districtCode string) string {
	uuid := uuid.New()
	return fmt.Sprintf("MQCT_%s_%s", districtCode, strings.Replace(uuid.String(), "-", "", -1))
}

func GetTransactionPrepareQueueName(topic string) string {
	return fmt.Sprintf("%s%s%s", REDIS_MQ_TRANSACTION_PRE_QUEUE_KEY, topic, NAME_VERSION)
}

func GetDeathQueueName() string {
	return fmt.Sprintf("%s%sdeath_message%s", REDIS_MQ_DEATH_QUEUE_KEY, STREAM_NAME, NAME_VERSION)
}

func GetMessageKey(topic string, tag string) string {
	return fmt.Sprintf("%s_%s", topic, tag)
}

func getTransactionDeathQueueName() string {
	return fmt.Sprintf("%s%s", REDIS_MQ_TRANSACTION_DEATH_QUEUE_KEY, NAME_VERSION)
}

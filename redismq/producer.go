package redismq

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/redis/go-redis/v9"
	"strings"
	"unibee/utility"
)

func Send(message *Message) (bool, error) {
	return sendMessage(message, "ProducerWrapper")
}
func SendTransaction(message *Message, transactionExecuter func(messageToSend *Message) (TransactionStatus, error)) (bool, error) {
	if strings.Compare(message.Tag, "blank") == 0 {
		return false, errors.New("blank tag message")
	}

	if message.StartDeliverTime > 0 {
		return false, errors.New("delay message not support transaction")
	}

	send, err := sendTransactionPrepareMessage(message)
	if err != nil || !send {
		return send, err
	}
	status, err := transactionExecuter(message)
	if status == RollbackTransaction {
		//事务执行失败，回滚半消息
		_, rollBackErr := rollbackTransactionPrepareMessage(message)
		if rollBackErr != nil {
			fmt.Printf("rollbackTransactionPrepareMessage err:%s rollBackError:%s\n", err, rollBackErr)
		}
		return false, err
	} else if status == CommitTransaction {
		//事务执行成功，提交半消息，如提交失败，需使用 实现相应Checker 保障消息一致性 todo mark
		return commitTransactionPrepareMessage(message)
	} else {
		//未知状态，一般在用户无法确定事务是成功还是失败时使用，对于未知状态的事务，服务端会定期进行事务回查
		return false, errors.New("unknown transaction status")
	}
}

func sendDelayMessage(message *Message) bool {
	utility.Assert(message.StartDeliverTime-gtime.Now().Timestamp() > 0, "StartDeliverTime Invalid, should > now")
	send, err := SendDelay(message, message.StartDeliverTime-gtime.Now().Timestamp())
	fmt.Printf("Redismq SendDelayMessage result:%v", send)
	if err != nil {
		return false
	}
	return send
}

func sendMessage(message *Message, source string) (bool, error) {
	if strings.Compare(message.Tag, "blank") == 0 {
		return false, errors.New("blank空消息")
	}

	message.SendTime = utility.CurrentTimeMillis()
	utility.Assert(len(message.MessageId) == 0, "Send Stream Need Blank MessageId")
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// Close Conn
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("sendMessage error:%s\n", err)
		}
	}(client)

	// 发送消息到 Stream
	streamMessageId, err := client.XAdd(context.Background(), message.toStreamAddArgsValues(GetQueueName(message.Topic))).Result()
	if err != nil {
		return false, errors.New(fmt.Sprintf("MQ STREAM Send MQStream exception:%s queueName=%s message:%v\n", err, GetQueueName(message.Topic), utility.MarshalToJsonString(message)))
	}
	message.MessageId = streamMessageId
	fmt.Printf("MQ STREAM Send Stream Success, Source:%s QueueName=%s MessageId=%v\n", source, GetQueueName(message.Topic), message.MessageId)
	return true, nil
}

func sendTransactionPrepareMessage(message *Message) (bool, error) {
	if strings.Compare(message.Tag, "blank") == 0 {
		return false, errors.New("Blank Message")
	}
	message.MessageId = GenerateUniqueNo(message.Topic)
	message.SendTime = utility.CurrentTimeMillis()
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// Close Conn
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("sendTransactionPrepareMessage error:%s\n", err)
		}
	}(client)
	messageJson, err := gjson.Marshal(message)

	jsonString := string(messageJson)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Send MQ Transaction Pre exception:%s message:%v\n", err.Error(), message))
	}
	// 执行事务
	_, err = client.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		// 在事务中执行多个命令
		//pipe.Incr(context.Background(), key)  // 递增键的值
		//pipe.Expire(context.Background(), key, 10*time.Second)  // 设置键的过期时间
		pipe.Set(context.Background(), message.MessageId, jsonString, -1)
		pipe.LPush(context.Background(), GetTransactionPrepareQueueName(message.Topic), message.MessageId)
		return nil
	})

	if err != nil {
		return false, errors.New(fmt.Sprintf("Send MQ Transaction Pre  exception:%s message:%v\n", err.Error(), message))
	}
	return true, nil
}

func rollbackTransactionPrepareMessage(message *Message) (bool, error) {
	return delTransactionPrepareMessage(message)
}

func delTransactionPrepareMessage(message *Message) (bool, error) {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// Close Conn
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("delTransactionPrepareMessage error:%s\n", err)
		}
	}(client)

	// 执行事务
	_, err := client.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		// 在事务中执行多个命令
		pipe.Del(context.Background(), message.MessageId)
		pipe.LRem(context.Background(), GetTransactionPrepareQueueName(message.Topic), 1, message.MessageId)
		return nil
	})

	if err != nil {
		return false, errors.New(fmt.Sprintf("Del MQ Transaction Pre  exception:%s message:%v\n", err, message))
	}
	fmt.Printf("rollbackTransactionPrepareMessage message:%v\n", message)
	return true, nil
}

func commitTransactionPrepareMessage(message *Message) (bool, error) {
	oldMessageId := message.MessageId
	message.MessageId = ""
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// Close Conn
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("commmitTransactionPrepareMessage error:%s\n", err)
		}
	}(client)
	streamMessageId := ""
	// 执行事务提交半消息到 Stream
	_, err := client.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		// 在事务中执行多个命令
		// 发送 Stream 消息
		streamMessageId, _ = client.XAdd(context.Background(), message.toStreamAddArgsValues(GetQueueName(message.Topic))).Result()
		message.MessageId = streamMessageId
		// 删除事务半消息
		pipe.Del(context.Background(), oldMessageId)
		pipe.LRem(context.Background(), GetTransactionPrepareQueueName(message.Topic), 1, oldMessageId)
		return nil
	})

	if err != nil {
		return false, errors.New(fmt.Sprintf("Commit MQ Transaction Pre  exception:%s message:%v\n", err, message))
	}
	fmt.Printf("Redismq commitTransactionPrepareMessage success message:%v prepareMessageId:%s targetMessageId:%s ", message, oldMessageId, streamMessageId)
	return true, nil
}

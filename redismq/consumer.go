package redismq

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/redis/go-redis/v9"
	"go-oversea-pay/utility"
	"net"
	"strings"
	"time"
)

const (
	GroupId = "GID_UniBee_Recurring"
)

var consumerName = ""

func StartRedisMqConsumer() {
	innerSettingConsumerName()
	if len(consumerName) == 0 {
		fmt.Println("StartRedisMqConsumer Failed While ConsumerName Invalid")
		return
	}
	StartDelayBackgroundThread()
	fmt.Println("Redismq Start Delay Queue！！！！！！")
	deathQueueName := GetDeathQueueName()
	createStreamGroup(deathQueueName, "death_message")
	fmt.Printf("Redismq Stream Init Death Queue deathQueueName:%s", deathQueueName)
	innerLoadTransactionChecker()
	fmt.Println("Redismq Finish Transaction Check Loader ！！！！！！")
	innerLoadConsumer()
	fmt.Println("Redismq Finish Default MQ Subscribe！！！！！！")
	startScheduleTrimStream()
	fmt.Println("Redismq Finish Queue Length Cut！！！！！！")
}

func innerSettingConsumerName() {
	// 获取本机的所有网络接口信息
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 遍历所有网络接口
	for _, iface := range interfaces {
		// 排除 lo（loopback）接口
		if iface.Flags&net.FlagLoopback == 0 {
			// 获取该网络接口的所有地址信息
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			// 遍历所有地址
			for _, addr := range addrs {
				// 将地址转换为 IP 地址
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				// 判断是否为 IPv4 地址
				if ip.To4() != nil {
					fmt.Printf("IPv4 Address: %s\n", ip)
					consumerName = ip.String()
				}
			}
		}
	}
}

func createStreamGroup(queueName string, topic string) {
	tryCreateGroup(queueName, topic)
	tryCreateConsumer(queueName, topic)
}

func tryCreateGroup(queueName string, topic string) {
	defer func() {
		if exception := recover(); exception != nil {
			fmt.Printf("Redismq Stream Init TryCreateGroup panic error:%s\n", exception)
			return
		}
	}()
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("sendMessage error:%s\n", err)
		}
	}(client)
	message := &Message{
		Topic: topic,
		Tag:   "blank",
		Body:  "test",
	}
	// 发送一条测试消息到 Stream
	_, err := client.XAdd(context.Background(), message.toStreamAddArgsValues(queueName)).Result()
	if err != nil {
		fmt.Printf("MQ STREAM初始化 Group Failure Or Group Exsit exception:%s queueName:%s group:%s\n", err, queueName, GroupId)
	}
	found := false
	groups, _ := client.XInfoGroups(context.Background(), queueName).Result()
	for _, group := range groups {
		if group.Name == GroupId {
			found = true
		}
	}
	if !found {
		//尝试创建 Group
		// 创建消费者组
		if err := client.XGroupCreateMkStream(context.Background(), queueName, GroupId, "$").Err(); err != nil {
			fmt.Printf("MQ STREAM GroupId exsit queueName:%s groupId:%s err:%s", queueName, GroupId, err)
			return
		} else {
			fmt.Printf("MQ STREAM init queueName:%s groupId:%s ", queueName, GroupId)
		}
	}
}

func tryCreateConsumer(queueName string, topic string) {
	defer func() {
		if exception := recover(); exception != nil {
			fmt.Printf("Redismq Stream init queue tryCreateConsumer panic error:%s\n", exception)
			return
		}
	}()
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("sendMessage error:%s\n", err)
		}
	}(client)
	if _, err := client.XGroupCreateConsumer(context.Background(), queueName, GroupId, consumerName).Result(); err != nil {
		fmt.Printf("MQ STREAM consumerName failure or consumerName exsit queueName:%s groupId:%s consumerName:%s err:%s", queueName, GroupId, consumerName, err)
	} else {
		fmt.Printf("MQ STREAM init queueName:%s groupId:%s consumerName:%s", queueName, GroupId, consumerName)
	}
}

func innerLoadConsumer() {
	for _, topic := range Topics {
		blockConsumerTopic(topic)
	}
}

func innerLoadTransactionChecker() {
	//checkers := checker.Checkers
	// Deprecated
}

func blockConsumerTopic(topic string) {
	createStreamGroup(GetQueueName(topic), topic)
	createStreamGroup(getBackupQueueName(topic), topic)
	// start background
	go loopConsumer(topic)
	go loopTransactionChecker(topic)
}

func loopConsumer(topic string) {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Closs Redis Stream Client error:%s\n", err)
		}
	}(client)
	for {
		var err error
		defer func() {
			if exception := recover(); exception != nil {
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
				}
				fmt.Printf("MQ STREAM Redismq Stream loopConsumer Redis Error topic:%s panic error:%s\n", topic, err.Error())
				return
			}
		}()
		count := 0
		message := blockReceiveConsumerMessage(client, topic)
		if message != nil {
			if consumer := getConsumer(message); consumer != nil {
				runConsumeMessage(consumer, message)
				//找不到本地订阅方,塞回队列，ToDo mark 实现group拉取方式应该丢弃
			} else {
				fmt.Printf("MQ STREAM Redismq Stream Receive Group:{} No Comsumer Drop message::%v\n", message)
				messageAck(message)
			}
			count++
		}
		//当所有的队列都为空时休眠1s
		if count == len(Topics) {
			time.Sleep(1 * time.Second)
		}
	}
}

func loopTransactionChecker(topic string) {
	for {
		defer func() {
			if exception := recover(); exception != nil {
				fmt.Printf("Redismq Stream MQ Query Transaction Pre Redis Error loopTransactionChecker topic:%s panic error:%s\n", topic, exception)
				return
			}
		}()
		messages := fetchTransactionPrepareMessagesForChecker(topic)
		for _, message := range messages {
			if ck := Checkers()[GetMessageKey(message.Topic, message.Tag)]; ck != nil {
				status := ck.Checker(message)
				if status == CommitTransaction {
					_, _ = commitTransactionPrepareMessage(message)
				} else if status == RollbackTransaction {
					_, _ = rollbackTransactionPrepareMessage(message)
				} else {
					//保存数据次数，最多50次重试
					if (utility.CurrentTimeMillis() - message.SendTime) > 1000*60*60*8 {
						//超过8小时，事务消息丢入死信队列
						putMessageToTransactionDeathQueue(topic, message)
					}
				}
			} else {
				// todo mark 检查优化处理没有checker 的情况，超过一定时间删除半消息
				if (utility.CurrentTimeMillis() - message.SendTime) > 1000*60*60*24*7 {
					//超过7天，事务消息回滚
					_, _ = rollbackTransactionPrepareMessage(message)
				}
			}
			time.Sleep(1 * time.Second)
		}
		time.Sleep(60 * time.Second)
	}
}

func getConsumer(message *Message) IMessageListener {
	if strings.Compare(message.Tag, "blank") == 0 {
		return nil
	}
	return Listeners()[GetMessageKey(message.Topic, message.Tag)]
}

func runConsumeMessage(consumer IMessageListener, message *Message) {
	var err error
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			fmt.Printf("MQ STREAM Redismq runConsumeMessage panic error:%s\n", err.Error())
			return
		}
	}()
	if message.isBoardCastingMessage() {
		// todo mark 出现这种情况属于bug
		fmt.Printf("RedisMQ_Receive Stream Message Exception Group Receive Boardcast，Drop messageKey:%s message:%v\n", GetMessageKey(message.Topic, message.Tag), message)
		return
	}
	cost := utility.CurrentTimeMillis()
	if message.SendTime > 0 {
		cost = utility.CurrentTimeMillis() - message.SendTime
		//历史消息没有过期时间
		if (utility.CurrentTimeMillis() - message.SendTime) > 1000*60*60*24*3 {
			//超过3天消息定义为过期消息，丢弃
			fmt.Printf("RedisMQ_Receive Stream Message Exception After 3 Days Drop Expired messageKey:%s message:%v\n ", GetMessageKey(message.Topic, message.Tag), message)
			return
		}
	} else {
		cost = 0
	}
	go func() {
		ctx := context.Background()
		defer func() {
			if exception := recover(); exception != nil {
				fmt.Printf("RedisMQ_Receive Stream Message Error message:%v panic error:%s\n", message, exception)
				if pushTaskToResumeLater(consumer, message) {
					messageAck(message)
				} else {
					// todo mark 进入Resumer流程失败,防止消息丢失
				}
				return
			}
		}()
		time.Sleep(2 * time.Second)
		action := consumer.Consume(ctx, message)
		fmt.Printf("RedisMQ_Receive Stream Message messageKey:%s result:%d message:%v cost:%dms", GetMessageKey(message.Topic, message.Tag), action, message, cost)
		if action == ReconsumeLater {
			if pushTaskToResumeLater(consumer, message) {
				messageAck(message)
			} else {
				// todo mark 进入Resumer流程失败,防止消息丢失
			}
		} else {
			messageAck(message)
		}
	}()
}

func messageAck(message *Message) {
	var err error
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			fmt.Printf("MQ STREAM Redismq MessageAck panic error:%s\n", err.Error())
			return
		}
	}()
	//todo mark messageId java 中有特殊处理
	//if strings.Contains(message.MessageId, "-") {
	//
	//} else {
	//
	//}
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("sendMessage error:%s\n", err)
		}
	}(client)
	streamName := GetQueueName(message.Topic)
	ackResult, err := client.XAck(context.Background(), streamName, GroupId, message.MessageId).Result()
	if err != nil {
		fmt.Printf("MQStream ack message:%v panic error:%s\n", message, err)
		return
	}
	fmt.Printf("MQStream ack streamMessageId:%s streamName:%s ackResult:%d", message.MessageId, streamName, ackResult)
}

func blockReceiveConsumerMessage(client *redis.Client, topic string) *Message {
	var err error
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception)
			}
			fmt.Printf("MQ STREAM blockReceiveConsumerMessage topic:%s panic error:%v %v\n", topic, err.Error(), exception)
			return
		}
	}()
	streamName := GetQueueName(topic)
	//fmt.Printf("MQStream XReadGroup blockReceiveConsumerMessage streamName=%s\n", streamName)
	result, err := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    GroupId,
		Consumer: consumerName,
		Streams:  []string{streamName, ">"},
		Count:    1,
		Block:    60 * time.Second,
		NoAck:    true,
	}).Result()
	if err != nil {
		fmt.Printf("MQStream blockReceiveConsumerMessage streamName=%s err=%s\n", streamName, err.Error())
		return nil
	}
	if len(result) == 1 && len(result[0].Messages) == 1 {
		//只获取一个
		messageId := result[0].Messages[0].ID
		value := result[0].Messages[0].Values
		message := Message{}
		message.MessageId = messageId
		message.getUniqueKey()
		message.paseStreamMessage(value)
		return &message
	}
	return nil
}

func pushTaskToResumeLater(consumer IMessageListener, message *Message) bool {
	ResumeTimesMax := 10
	if message.ReconsumeTimes > ResumeTimesMax {
		return putMessageToDeathQueue(message.Topic, message.MessageId, message)
	} else {
		//延长时间重新发送消息
		message.ReconsumeTimes = message.ReconsumeTimes + 1
		message.StartDeliverTime = utility.CurrentTimeMillis() + (20 * 1000 * int64(message.ReconsumeTimes)) // 20 秒
		return sendDelayMessage(message)
	}
}

func putMessageToDeathQueue(topic string, id string, message *Message) bool {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("sendMessage error:%s\n", err)
		}
	}(client)
	streamMessageId, err := client.XAdd(context.Background(), message.toStreamAddArgsValues(GetDeathQueueName())).Result()
	if err != nil {
		fmt.Printf("MQStream push message to death exception:%s messageId:%s", err, message.MessageId)
		return false
	}
	fmt.Printf("MQ push message to death, message=%v deathMessageId:%s", message, streamMessageId)
	return true
}

func putMessageToTransactionDeathQueue(topic string, message *Message) bool {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("MQ push transaction message to death error:%s\n", err)
		}
	}(client)

	// 执行事务
	_, err := client.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
		// 在事务中执行多个命令
		//pipe.Incr(context.Background(), key)  // 递增键的值
		//pipe.Expire(context.Background(), key, 10*time.Second)  // 设置键的过期时间
		pipe.LRem(context.Background(), GetTransactionPrepareQueueName(topic), 1, message.MessageId)
		pipe.RPush(context.Background(), getTransactionDeathQueueName(), message.MessageId)
		return nil
	})

	if err != nil {
		fmt.Printf("transaction message to death and delete exception:%s message:%v\n", err, message)
		return false
	}
	return true
}

func fetchTransactionPrepareMessagesForChecker(topic string) []*Message {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("MQ redis error:%s\n", err)
		}
	}(client)

	result, err := client.LRange(context.Background(), GetTransactionPrepareQueueName(topic), 0, -1).Result()
	if err != nil {
		return []*Message{}
	}
	var messages = make([]*Message, 0)
	for _, messageId := range result {
		if len(messageId) > 0 {
			value, _ := client.Get(context.Background(), messageId).Result()
			if len(value) > 0 {
				var message *Message
				err = gjson.Unmarshal([]byte(value), &message) // Unmarshal todo mark 加上 &
				if err == nil {
					messages = append(messages, message)
				}
			} else {
				fmt.Printf("transaction pre message messageId:%s\n", messageId)
			}
		}
	}
	return messages
}

func startScheduleTrimStream() {
	var maxLen = 10000
	go func() {
		client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
		// 关闭连接
		defer func(client *redis.Client) {
			err := client.Close()
			if err != nil {
				fmt.Printf("MQ redis error:%s\n", err)
			}
		}(client)
		for {
			defer func() {
				if exception := recover(); exception != nil {
					fmt.Printf("startScheduleTrimStream exception:%s\n", exception)
					return
				}
			}()
			for _, topic := range Topics {
				queueName := GetQueueName(topic)
				client.XTrimMaxLen(context.Background(), queueName, int64(maxLen))
				fmt.Printf("MQ STREAM Cut maxLen:%d queueName:%s group:%s consumerName:%s\n", maxLen, queueName, GroupId, consumerName)
				consumersCheck(queueName)
				queueName = getBackupQueueName(topic)
				client.XTrimMaxLen(context.Background(), queueName, int64(maxLen))
				fmt.Printf("MQ STREAM Cut maxLen:%d queueName:%s group:%s consumerName:%s\n", maxLen, queueName, GroupId, consumerName)
				consumersCheck(queueName)
			}

			queueName := GetDeathQueueName()
			client.XTrimMaxLen(context.Background(), queueName, int64(maxLen))
			fmt.Printf("MQ STREAM Cut maxLen:%d queueName:%s group:%s consumerName:%s\n", maxLen, queueName, GroupId, consumerName)
			consumersCheck(queueName)

			time.Sleep(1000 * 60 * 10 * time.Second) //10分钟修剪一次
		}
	}()
}

func consumersCheck(queueName string) {
	//todo mark
}

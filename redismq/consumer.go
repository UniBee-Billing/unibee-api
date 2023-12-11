package redismq

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/redis/go-redis/v9"
	"go-oversea-pay/redismq/checker"
	"go-oversea-pay/redismq/listener"
	"go-oversea-pay/utility"
	"net"
	"strings"
	"time"
)

const (
	GroupId = "GID_hk_golang_overseapay"
)

var consumerName = ""

func StartRedisMqConsumer() {
	innerSettingConsumerName()
	if len(consumerName) == 0 {
		fmt.Println("StartRedisMqConsumer failed while ConsumerName invalid")
		return
	}
	StartDelayBackgroundThread()
	fmt.Println("redismq 启动延迟队列后台任务！！！！！！")
	deathQueueName := GetDeathQueueName()
	createStreamGroup(deathQueueName, "death_message")
	fmt.Printf("redismq Stream 初始化死信队列 deathQueueName:%s", deathQueueName)
	innerLoadTransactionChecker()
	fmt.Println("redismq 事务checkloader汇总完毕！！！！！！")
	innerLoadConsumer()
	fmt.Println("redismq 默认订阅信息汇总完毕！！！！！！")
	startScheduleTrimStream()
	fmt.Println("redismq 定时修剪长度任务完毕！！！！！！")
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
					consumerName = string(ip)
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
			fmt.Printf("redismq Stream 初始化队列 tryCreateGroup panic error:%s\n", exception)
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
		Body:  []byte("test"),
	}
	// 发送一条测试消息到 Stream
	_, err := client.XAdd(context.Background(), message.toStreamAddArgsValues(queueName)).Result()
	if err != nil {
		fmt.Printf("MQ STREAM初始化 Group失败或Group已存在 exception:%s queueName:%s group:%s\n", err, queueName, GroupId)
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
			fmt.Printf("MQ STREAM 已存在GroupId queueName:%s groupId:%s err:%s", queueName, GroupId, err)
			return
		} else {
			fmt.Printf("MQ STREAM初始化 queueName:%s groupId:%s ", queueName, GroupId)
		}
	}
}

func tryCreateConsumer(queueName string, topic string) {
	defer func() {
		if exception := recover(); exception != nil {
			fmt.Printf("redismq Stream 初始化队列 tryCreateConsumer panic error:%s\n", exception)
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
		fmt.Printf("MQ STREAM consumerName失败或consumerName已存在 queueName:%s groupId:%s consumerName:%s err:%s", queueName, GroupId, consumerName, err)
	} else {
		fmt.Printf("MQ STREAM初始化 queueName:%s groupId:%s consumerName:%s", queueName, GroupId, consumerName)
	}
}

func innerLoadConsumer() {
	for _, topic := range listener.Topics {
		blockConsumerTopic(topic)
	}
}

func innerLoadTransactionChecker() {
	//checkers := checker.Checkers
	// 废弃，不需要了
}

func blockConsumerTopic(topic string) {
	createStreamGroup(GetQueueName(topic), topic)
	createStreamGroup(getBackupQueueName(topic), topic)
	//启动协程
	go loopConsumer(topic)
	go loopTransactionChecker(topic)
}

func loopConsumer(topic string) {
	for {
		defer func() {
			if exception := recover(); exception != nil {
				fmt.Printf("redismq Stream loopConsumer 消息获取Redis异常 topic:%s panic error:%s\n", topic, exception)
				return
			}
		}()
		count := 0
		message := blockReceiveConsumerMessage(topic)
		if message != nil {
			if consumer := getConsumer(message); consumer != nil {
				runConsumeMessage(consumer, message)
				//找不到本地订阅方,塞回队列，ToDo mark 实现group拉取方式应该丢弃
			} else {
				fmt.Printf("redismq Stream 收到当前Group:{}无消费者消息 丢弃消息 message::%v\n", message)
				messageAck(message)
			}
			count++
		}
		//当所有的队列都为空时休眠1s
		if count == len(listener.Topics) {
			time.Sleep(1 * time.Second)
		}
	}
}

func loopTransactionChecker(topic string) {
	for {
		defer func() {
			if exception := recover(); exception != nil {
				fmt.Printf("redismq Stream MQ事务半消息获取Redis异常 loopTransactionChecker topic:%s panic error:%s\n", topic, exception)
				return
			}
		}()
		messages := fetchTransactionPrepareMessagesForChecker(topic)
		for _, message := range messages {
			if ck := checker.Checkers()[GetMessageKey(message.Topic, message.Tag)]; ck != nil {
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

func getConsumer(message *Message) listener.IMessageListener {
	if strings.Compare(message.Tag, "blank") == 0 {
		return nil
	}
	return listener.Listeners()[GetMessageKey(message.Topic, message.Tag)]
}

func runConsumeMessage(consumer listener.IMessageListener, message *Message) {
	if message.isBoardCastingMessage() {
		// todo mark 出现这种情况属于bug
		fmt.Printf("RedisMQ_收到Stream消息 Exception Group通道收到广播消息，丢弃 messageKey:%s message:%v\n", GetMessageKey(message.Topic, message.Tag), message)
		return
	}
	cost := utility.CurrentTimeMillis()
	if message.SendTime > 0 {
		cost = utility.CurrentTimeMillis() - message.SendTime
		//历史消息没有过期时间
		if (utility.CurrentTimeMillis() - message.SendTime) > 1000*60*60*24*3 {
			//超过3天消息定义为过期消息，丢弃
			fmt.Printf("RedisMQ_收到Stream消息 Exception 消息超过3天 过期丢弃 messageKey:%s message:%v\n ", GetMessageKey(message.Topic, message.Tag), message)
			return
		}
	} else {
		cost = 0
	}
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				fmt.Printf("RedisMQ_收到Stream消息 执行消费任务异常 message:%v panic error:%s\n", message, exception)
				if pushTaskToResumeLater(consumer, message) {
					messageAck(message)
				} else {
					// todo mark 进入Resumer流程失败,防止消息丢失
				}
				return
			}
		}()
		time.Sleep(2 * time.Second)
		action := consumer.Consume(message)
		fmt.Printf("RedisMQ_收到Stream消息 messageKey:%s 执行结果:%d message:%v cost:%dms", GetMessageKey(message.Topic, message.Tag), action, message, cost)
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
	defer func() {
		if exception := recover(); exception != nil {
			fmt.Printf("MQStream消息发送ACK异常 message:%v panic error:%s\n", message, exception)
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
		fmt.Printf("MQStream消息发送ACK异常 message:%v panic error:%s\n", message, err)
		return
	}
	fmt.Printf("MQStream消息发送ACK streamMessageId:%s streamName:%s ackResult:%d", message.MessageId, streamName, ackResult)
}

func blockReceiveConsumerMessage(topic string) *Message {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("sendMessage error:%s\n", err)
		}
	}(client)
	streamName := GetQueueName(topic)
	result, err := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    GroupId,
		Consumer: consumerName,
		Streams:  []string{streamName},
		Count:    1,
		Block:    60 * time.Second,
		NoAck:    true,
	}).Result()
	if err != nil {
		fmt.Printf("MQStream消息获取Redis异常 exception=%s", err)
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

func pushTaskToResumeLater(consumer listener.IMessageListener, message *Message) bool {
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
		fmt.Printf("MQStream消息推入死信队列异常 exception:%s messageId:%s", err, message.MessageId)
		return false
	}
	fmt.Printf("MQ消息推入死信队列, message=%v deathMessageId:%s", message, streamMessageId)
	return true
}

func putMessageToTransactionDeathQueue(topic string, message *Message) bool {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("MQ事务消息推入死信队列异常 error:%s\n", err)
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
		fmt.Printf("事务消息推入死信队列删除不成功 exception:%s message:%v\n", err, message)
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
			fmt.Printf("MQ消息获取Redis异常 error:%s\n", err)
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
				err = gjson.Unmarshal([]byte(value), message)
				if err == nil {
					messages = append(messages, message)
				}
			} else {
				fmt.Printf("事务半消息体未找到 messageId:%s\n", messageId)
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
				fmt.Printf("MQ消息获取Redis异常 error:%s\n", err)
			}
		}(client)
		for {
			defer func() {
				if exception := recover(); exception != nil {
					fmt.Printf("startScheduleTrimStream exception:%s\n", exception)
					return
				}
			}()
			for _, topic := range listener.Topics {
				queueName := GetQueueName(topic)
				client.XTrimMaxLen(context.Background(), queueName, int64(maxLen))
				fmt.Printf("MQ STREAM 修剪长度 maxLen:%d queueName:%s group:%s consumerName:%s\n", maxLen, queueName, GroupId, consumerName)
				consumersCheck(queueName)
				queueName = getBackupQueueName(topic)
				client.XTrimMaxLen(context.Background(), queueName, int64(maxLen))
				fmt.Printf("MQ STREAM 修剪长度 maxLen:%d queueName:%s group:%s consumerName:%s\n", maxLen, queueName, GroupId, consumerName)
				consumersCheck(queueName)
			}

			queueName := GetDeathQueueName()
			client.XTrimMaxLen(context.Background(), queueName, int64(maxLen))
			fmt.Printf("MQ STREAM 修剪长度 maxLen:%d queueName:%s group:%s consumerName:%s\n", maxLen, queueName, GroupId, consumerName)
			consumersCheck(queueName)

			time.Sleep(1000 * 60 * 10 * time.Second) //10分钟修剪一次
		}
	}()
}

func consumersCheck(queueName string) {
	//todo mark
}

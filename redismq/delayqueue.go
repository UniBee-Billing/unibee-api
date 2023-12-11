package redismq

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/redis/go-redis/v9"
	"go-oversea-pay/utility"
	"strconv"
	"time"
)

const (
	MQ_DELAY_QUEUE_NAME = "MQ_DELAY_QUEUE_SET"
)

func StartDelayBackgroundThread() {
	go func() {
		for {
			defer func() {
				if exception := recover(); exception != nil {
					fmt.Printf("redismq polligCore 异常 panic error:%s\n", exception)
					return
				}
			}()
			polling()
			time.Sleep(10 * time.Second)
		}
	}()
}

func polling() {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("redismq run error:%s\n", err)
		}
	}(client)

	result, err := client.Keys(context.Background(), MQ_DELAY_QUEUE_NAME).Result()
	if err != nil {
		return
	}
	for _, key := range result {
		pollingCore(key)
	}
}

func pollingCore(key string) {
	defer func() {
		if exception := recover(); exception != nil {
			fmt.Printf("redismq polligCore 异常 panic error:%s\n", exception)
			return
		}
	}()

	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("redismq polligCore error:%s\n", err)
		}
	}(client)
	result, err := client.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(utility.CurrentTimeMillis(), 10), // todo mark formatint base 值可能存在问题
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		return
	}
	if len(result) == 0 {
		fmt.Printf("redismq 延时队列[%s]没有过期节点\n", MQ_DELAY_QUEUE_NAME)
		return
	}
	for _, messageJson := range result {
		fmt.Printf("redismq 从延时队列移除节点[%s]\n", messageJson)
		var message *Message

		// 使用 gjson.Unmarshal 将 JSON 字符串解析成结构体
		err := gjson.Unmarshal([]byte(messageJson), message)
		if err != nil {
			fmt.Printf("Error:%v\n", err)
			return
		}

		message.StartDeliverTime = 0
		message.MessageId = ""
		send, err := sendMessage(message, "DelayQueue")
		fmt.Printf("redismq 从延时队列移除节点,并发送消息[%v], error:[%s]\n", send, err)
	}
}

func SendDelay(message *Message, delay int64) (bool, error) {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("redismq SendDelay error:%s\n", err)
		}
	}(client)
	messageJson, err := gjson.Marshal(message)
	if err != nil {
		return false, errors.New(fmt.Sprintf("发送MQ事务半消息异常 exception:%s message:%v\n", err, message))
	}
	jsonString := string(messageJson)
	score := utility.CurrentTimeMillis() + delay
	result, err := client.ZAdd(context.Background(), MQ_DELAY_QUEUE_NAME, redis.Z{
		Score:  float64(score),
		Member: jsonString,
	}).Result()
	if err != nil {
		return false, errors.New(fmt.Sprintf("发送MQ事务半消息异常 exception:%s message:%v\n", err, message))
	}
	fmt.Printf("redismq 向延时队列投放任务,队列名称[%s],任务[%s],score[%d],result[%v]\n", MQ_DELAY_QUEUE_NAME, messageJson, score, result)
	return true, nil
}

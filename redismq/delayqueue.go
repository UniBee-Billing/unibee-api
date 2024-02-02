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
					fmt.Printf("Redismq polligCore panic error:%s\n", exception)
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
			fmt.Printf("Redismq run error:%s\n", err)
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
			fmt.Printf("Redismq polligCore panic error:%s\n", exception)
			return
		}
	}()

	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Redismq polligCore error:%s\n", err)
		}
	}(client)
	result, err := client.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(utility.CurrentTimeMillis(), 10),
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		return
	}
	if len(result) == 0 {
		fmt.Printf("Redismq Delay Queue[%s] No Queue\n", MQ_DELAY_QUEUE_NAME)
		return
	}
	for _, messageJson := range result {
		fmt.Printf("Redismq Delete From Delay Queue[%s]\n", messageJson)
		var message *Message

		// 使用 gjson.Unmarshal 将 JSON 字符串解析成结构体
		err := gjson.Unmarshal([]byte(messageJson), &message) // Unmarshal todo mark 加上 &
		if err != nil {
			fmt.Printf("Error:%v\n", err)
			return
		}

		message.StartDeliverTime = 0
		message.MessageId = ""
		send, err := sendMessage(message, "DelayQueue")
		fmt.Printf("Redismq Delete From Delay Queue,And Send[%v], error:[%s]\n", send, err)
	}
}

func SendDelay(message *Message, delay int64) (bool, error) {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Redismq SendDelay error:%s\n", err)
		}
	}(client)
	messageJson, err := gjson.Marshal(message)
	if err != nil {
		return false, errors.New(fmt.Sprintf("SendDelay exception:%s message:%v\n", err.Error(), message))
	}
	jsonString := string(messageJson)
	score := utility.CurrentTimeMillis() + delay
	result, err := client.ZAdd(context.Background(), MQ_DELAY_QUEUE_NAME, redis.Z{
		Score:  float64(score),
		Member: jsonString,
	}).Result()
	if err != nil {
		return false, errors.New(fmt.Sprintf("SendDelay exception:%s message:%v\n", err.Error(), message))
	}
	fmt.Printf("Redismq Push To Deplay Queue,Name[%s],Task[%s],Score[%d],Result[%v]\n", MQ_DELAY_QUEUE_NAME, messageJson, score, result)
	return true, nil
}

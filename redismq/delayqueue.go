package redismq

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/redis/go-redis/v9"
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
	ctx := context.Background()
	result, err := client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatInt(gtime.Now().Timestamp(), 10),
		Offset: 0,
		Count:  1,
	}).Result()
	if err != nil {
		return
	}
	if len(result) == 0 {
		g.Log().Debugf(ctx, "Redismq Delay Queue[%s] No Queue\n", MQ_DELAY_QUEUE_NAME)
		return
	}
	for _, messageJson := range result {
		var message *Message
		// 使用 gjson.Unmarshal 将 JSON 字符串解析成结构体
		err = gjson.Unmarshal([]byte(messageJson), &message)
		if err != nil {
			fmt.Printf("Redismq Unmarshal Message Error:[%v]\n", err)
			continue
		}
		err = client.ZRem(ctx, key, messageJson).Err()
		if err == nil {
			fmt.Printf("Redismq Delete From Delay Queue[%s]\n", messageJson)
			message.StartDeliverTime = 0
			message.MessageId = ""
			send, sendErr := sendMessage(message, "DelayQueue")
			fmt.Printf("Redismq Delete From Delay Queue,And Send[%v], Error:[%s] SendErr:[%s]\n", send, err, sendErr)
		} else {
			fmt.Printf("Redismq Delete From Delay Queue[%s] Err:[%s]\n", messageJson, err)
		}
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
	score := gtime.Now().Timestamp() + delay
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

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
	MqDelayQueueName = "MQ_DELAY_QUEUE_SET"
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
	// Close Conn
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("RedisMq run error:%s\n", err.Error())
		}
	}(client)

	result, err := client.Keys(context.Background(), MqDelayQueueName).Result()
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
			fmt.Printf("RedisMq polligCore panic error:%s\n", exception)
			return
		}
	}()

	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// Close Conn
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("RedisMq polligCore error:%s\n", err.Error())
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
		g.Log().Debugf(ctx, "RedisMq Delay Queue:%s No Queue\n", MqDelayQueueName)
		return
	}
	for _, messageJson := range result {
		var message *Message
		err = gjson.Unmarshal([]byte(messageJson), &message)
		if err != nil {
			fmt.Printf("RedisMq Unmarshal Message Error:[%v]\n", err)
			continue
		}
		err = client.ZRem(ctx, key, messageJson).Err()
		if err == nil {
			message.StartDeliverTime = 0
			message.MessageId = ""
			send, sendErr := sendMessage(message, "DelayQueue")
			if sendErr != nil {
				g.Log().Errorf(ctx, "RedisMq Delete From Delay Queue,And Send[%v], sendErr:%s\n", send, sendErr.Error())
			} else {
				g.Log().Debugf(ctx, "RedisMq Delete From Delay Queue,And Send:%v\n", send)
			}
		} else {
			g.Log().Debugf(ctx, "RedisMq Delete From Delay Err:%s\n", err.Error())
		}
	}
}

func SendDelay(message *Message, delay int64) (bool, error) {
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// Close Conn
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("RedisMq SendDelay error:%s \n", err.Error())
		}
	}(client)
	messageJson, err := gjson.Marshal(message)
	if err != nil {
		return false, errors.New(fmt.Sprintf("RedisMq SendDelay exception:%s message:%v\n", err.Error(), message))
	}
	jsonString := string(messageJson)
	score := gtime.Now().Timestamp() + delay
	_, err = client.ZAdd(context.Background(), MqDelayQueueName, redis.Z{
		Score:  float64(score),
		Member: jsonString,
	}).Result()
	if err != nil {
		return false, errors.New(fmt.Sprintf("SendDelay exception:%s message:%v\n", err.Error(), message))
	}
	//fmt.Printf("RedisMq Push To Deplay Queue,Name[%s],Task[%s],Score[%d],Result[%v]\n", MqDelayQueueName, messageJson, score, result)
	return true, nil
}

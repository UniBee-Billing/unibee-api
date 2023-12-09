package redismq

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

func testSend() {

	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())
	// 关闭连接
	defer func(client *redis.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	// Stream 名称
	streamName := "mystream"

	// 消息内容
	message := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	// 发送消息到 Stream
	_, err := client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: streamName,
		Values: message,
	}).Result()

	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Println("Message sent successfully!")

	// 等待一秒钟，确保消息被读取
	time.Sleep(time.Second)
}

func testReceive() {
	// 创建 Redis 客户端
	client := redis.NewClient(SharedConfig().GetRedisStreamConfig())

	// 关闭连接
	defer client.Close()

	// Stream 名称
	streamName := "mystream"

	// 消费者组
	consumerGroup := "mygroup"

	// 消费者名称
	consumerName := "myconsumer"

	// 创建消费者组
	client.XGroupCreateMkStream(context.Background(), streamName, consumerGroup, "$")

	// 循环消费消息
	for {
		// 从 Stream 中读取消息
		result, err := client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Block:    0,
			Count:    1,
		}).Result()

		if err != nil {
			fmt.Println("Error reading from stream:", err)
			return
		}

		// 处理消息
		for _, message := range result {
			for _, xMessage := range message.Messages {
				// 输出消息内容
				fmt.Printf("Received message: %v\n", xMessage.Values)
			}
		}

		// 暂停一秒钟，避免无限循环过快
		time.Sleep(time.Second)
	}
}

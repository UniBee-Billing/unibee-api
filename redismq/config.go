package redismq

import "github.com/redis/go-redis/v9"

type IRedisMqConConfig interface {
	GetRedisStreamConfig() (res *redis.Options)
}

var instance IRedisMqConConfig

func SharedConfig() IRedisMqConConfig {
	if instance == nil {
		panic("implement not found for interface IRedisMqConConfig, forgot register?")
	}
	return instance
}

func RegisterRedisMqConfig(i IRedisMqConConfig) {
	instance = i
}

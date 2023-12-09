package redismq

import (
	"github.com/redis/go-redis/v9"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/redismq"
)

type SRedisMqConfig struct{}

func (s SRedisMqConfig) GetRedisStreamConfig() (res *redis.Options) {
	return &redis.Options{
		Addr:     consts.GetNacosConfigInstance().RedisConfig.Address,
		Password: consts.GetNacosConfigInstance().RedisConfig.Pass,
		DB:       consts.GetNacosConfigInstance().RedisConfig.DB,
	}
}

func init() {
	redismq.RegisterRedisMqConfig(New())
}

func New() *SRedisMqConfig {
	return &SRedisMqConfig{}
}

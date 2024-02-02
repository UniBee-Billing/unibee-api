package redismq

import (
	"github.com/redis/go-redis/v9"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/redismq"
)

var (
	TopicTest1                      = redismq.MQTopicEnum{"unibee", "test1", "redismq测试1"}
	TopicTest2                      = redismq.MQTopicEnum{"unibee", "test2", "redismq测试2"}
	TopicChannelPayV2WebHookReceive = redismq.MQTopicEnum{"unibee", "channelpaywebhookreceivev2", "channelpaywebhookv2消息接收"}
	TopicPayCreated                 = redismq.MQTopicEnum{"unibee_pay", "paycreated", "支付单创建"}
	TopicPayCancelld                = redismq.MQTopicEnum{"unibee_pay", "payCancelld", "支付单取消成功"}
	TopicPayAuthorized              = redismq.MQTopicEnum{"unibee_pay", "payauthorized", "支付单授权"}
	TopicPaySuccess                 = redismq.MQTopicEnum{"unibee_pay", "paysuccess", "支付成功"}
	TopicRefundCreated              = redismq.MQTopicEnum{"unibee_refund", "refundcreated", "退款单创建"}
	TopicRefundSuccess              = redismq.MQTopicEnum{"unibee_refund", "refundsuccess", "退款成功"}
	TopicRefundFailed               = redismq.MQTopicEnum{"unibee_refund", "refundfailed", "退款失败"}
)

type SRedisMqConfig struct{}

func (s SRedisMqConfig) GetRedisStreamConfig() (res *redis.Options) {
	return &redis.Options{
		Addr:     consts.GetConfigInstance().RedisConfig.Default.Address,
		Password: consts.GetConfigInstance().RedisConfig.Default.Pass,
		DB:       consts.GetConfigInstance().RedisConfig.Default.DB,
	}
}

func init() {
	redismq.RegisterRedisMqConfig(New())
}

func New() *SRedisMqConfig {
	return &SRedisMqConfig{}
}

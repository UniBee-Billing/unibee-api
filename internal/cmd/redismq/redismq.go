package redismq

import (
	"github.com/redis/go-redis/v9"
	"go-oversea-pay/internal/consts"
	"go-oversea-pay/redismq"
)

var (
	TopicTest1                      = redismq.MQTopicEnum{"overseamq", "test1", "redismq测试1"}
	TopicTest2                      = redismq.MQTopicEnum{"overseamq", "test2", "redismq测试2"}
	TopicChannelPayV2WebHookReceive = redismq.MQTopicEnum{"overseamq", "channelpaywebhookreceivev2", "channelpaywebhookv2消息接收"}
	TopicPayCreated                 = redismq.MQTopicEnum{"overseamq_pay", "paycreated", "支付单创建"}
	TopicPayCancelld                = redismq.MQTopicEnum{"overseamq_pay", "payCancelld", "支付单取消成功"}
	TopicPayAuthorized              = redismq.MQTopicEnum{"overseamq_pay", "payauthorized", "支付单授权"}
	TopicPaySuccess                 = redismq.MQTopicEnum{"overseamq_pay", "paysuccess", "支付成功"}
	TopicRefundCreated              = redismq.MQTopicEnum{"overseamq_refund", "refundcreated", "退款单创建"}
	TopicRefundSuccess              = redismq.MQTopicEnum{"overseamq_refund", "refundsuccess", "退款成功"}
	TopicRefundFailed               = redismq.MQTopicEnum{"overseamq_refund", "refundfailed", "退款失败"}
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

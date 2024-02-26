package redismq

import (
	"github.com/redis/go-redis/v9"
	"unibee/internal/consts"
	"unibee/redismq"
)

var (
	TopicTest1                          = redismq.MQTopicEnum{"unibee", "test1", "redismq test1"}
	TopicTest2                          = redismq.MQTopicEnum{"unibee", "test2", "redismq test2"}
	TopicGatewayPayV2WebHookReceive     = redismq.MQTopicEnum{"unibee", "gatewaypaywebhookreceivev2", "gatewaypaywebhookv2"}
	TopicPayCreated                     = redismq.MQTopicEnum{"unibee_pay", "paycreated", "payment created"}
	TopicPayCancel                      = redismq.MQTopicEnum{"unibee_pay", "payCancelld", "payment cancelled"}
	TopicPayAuthorized                  = redismq.MQTopicEnum{"unibee_pay", "payauthorized", "payment authorized"}
	TopicPaySuccess                     = redismq.MQTopicEnum{"unibee_pay", "paysuccess", "payment success"}
	TopicRefundCreated                  = redismq.MQTopicEnum{"unibee_refund", "refundcreated", "refund created"}
	TopicRefundSuccess                  = redismq.MQTopicEnum{"unibee_refund", "refundsuccess", "refund success"}
	TopicRefundFailed                   = redismq.MQTopicEnum{"unibee_refund", "refundfailed", "refund success"}
	TopicSubscriptionCancel             = redismq.MQTopicEnum{"unibee_subscription", "subscription_cancelled", "subscription cancelled"}
	TopicSubscriptionExpire             = redismq.MQTopicEnum{"unibee_subscription", "subscription_expired", "subscription expired"}
	TopicSubscriptionIncomplete         = redismq.MQTopicEnum{"unibee_subscription", "subscription_incomplete", "subscription incomplete"}
	TopicSubscriptionCreate             = redismq.MQTopicEnum{"unibee_subscription", "subscription_created", "subscription created"}
	TopicSubscriptionCreatePaymentCheck = redismq.MQTopicEnum{"unibee_subscription", "subscription_create_payment_check", "subscription create payment check"}
	TopicSubscriptionPaymentSuccess     = redismq.MQTopicEnum{"unibee_subscription", "subscription_payment_success", "subscription payment success"}
	TopicMerchantWebhook                = redismq.MQTopicEnum{"unibee_merchant_webhook", "webhook", "merchant webhook"}
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

package redismq

import (
	"github.com/redis/go-redis/v9"
	"unibee/internal/cmd/config"
	"unibee/redismq"
)

var (
	TopicTest1                            = redismq.MQTopicEnum{"unibee", "test1", "redismq test1"}
	TopicTest2                            = redismq.MQTopicEnum{"unibee", "test2", "redismq test2"}
	TopicGatewayPayV2WebHookReceive       = redismq.MQTopicEnum{"unibee", "gatewaypaywebhookreceivev2", "gatewaypaywebhookv2"}
	TopicUserPaymentMethodChanged         = redismq.MQTopicEnum{"unibee_user", "payment_method_changed", "user change payment method"}
	TopicPaymentCreated                   = redismq.MQTopicEnum{"unibee_payment", "payment_created", "payment created"}
	TopicPaymentChecker                   = redismq.MQTopicEnum{"unibee_payment", "payment_checker", "payment status checker"}
	TopicPaymentCancel                    = redismq.MQTopicEnum{"unibee_payment", "payment_cancelled", "payment cancelled"}
	TopicPaymentAuthorized                = redismq.MQTopicEnum{"unibee_payment", "payment_authorized", "payment authorized"}
	TopicPaymentSuccess                   = redismq.MQTopicEnum{"unibee_payment", "payment_success", "payment success"}
	TopicRefundCreated                    = redismq.MQTopicEnum{"unibee_refund", "refund_created", "refund created"}
	TopicRefundChecker                    = redismq.MQTopicEnum{"unibee_refund", "refund_checker", "refund status checker"}
	TopicRefundSuccess                    = redismq.MQTopicEnum{"unibee_refund", "refund_success", "refund success"}
	TopicRefundFailed                     = redismq.MQTopicEnum{"unibee_refund", "refund_failed", "refund success"}
	TopicSubscriptionCancel               = redismq.MQTopicEnum{"unibee_subscription", "subscription_cancelled", "subscription cancelled"}
	TopicSubscriptionExpire               = redismq.MQTopicEnum{"unibee_subscription", "subscription_expired", "subscription expired"}
	TopicSubscriptionCreate               = redismq.MQTopicEnum{"unibee_subscription", "subscription_created", "subscription created"}
	TopicSubscriptionCreatePaymentCheck   = redismq.MQTopicEnum{"unibee_subscription", "subscription_create_payment_check", "subscription create payment check"}
	TopicSubscriptionIncomplete           = redismq.MQTopicEnum{"unibee_subscription", "subscription_incomplete", "subscription incomplete"}
	TopicSubscriptionPaymentSuccess       = redismq.MQTopicEnum{"unibee_subscription", "subscription_payment_success", "subscription payment success"}
	TopicSubscriptionAutoRenewSuccess     = redismq.MQTopicEnum{"unibee_subscription", "subscription_auto_renew_success", "subscription auto renew success"}
	TopicSubscriptionAutoRenewFailure     = redismq.MQTopicEnum{"unibee_subscription", "subscription_auto_renew_failure", "subscription auto renew failure"}
	TopicSubscriptionActiveWithoutPayment = redismq.MQTopicEnum{"unibee_subscription", "subscription_active_without_payment", "subscription become active without payment"}
	TopicMerchantWebhook                  = redismq.MQTopicEnum{"unibee_merchant_webhook", "webhook", "merchant webhook"}
	TopicInvoiceCreated                   = redismq.MQTopicEnum{"unibee_invoice", "invoice_created", "invoice created"}
	TopicInvoiceProcessed                 = redismq.MQTopicEnum{"unibee_invoice", "invoice_processed", "invoice processed"}
	TopicInvoicePaid                      = redismq.MQTopicEnum{"unibee_invoice", "invoice_paid", "invoice paid"}
	TopicInvoiceCancelled                 = redismq.MQTopicEnum{"unibee_invoice", "invoice_cancelled", "invoice cancelled"}
	TopicInvoiceFailed                    = redismq.MQTopicEnum{"unibee_invoice", "invoice_failed", "invoice failed"}
)

type SRedisMqConfig struct{}

func (s SRedisMqConfig) GetRedisStreamConfig() (res *redis.Options) {
	one := &redis.Options{
		Addr:     config.GetConfigInstance().RedisConfig.Default.Address,
		Password: config.GetConfigInstance().RedisConfig.Default.Pass,
		DB:       config.GetConfigInstance().RedisConfig.Default.DB,
	}
	return one
}

func init() {
	redismq.RegisterRedisMqConfig(New())
}

func New() *SRedisMqConfig {
	return &SRedisMqConfig{}
}

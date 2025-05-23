package redismq

import (
	redismq "github.com/jackyang-hk/go-redismq"
)

var (
	TopicTest1                            = redismq.MQTopicEnum{Topic: "unibee", Tag: "test1", Description: "redismq test1"}
	TopicUserPaymentMethodChange          = redismq.MQTopicEnum{Topic: "unibee_user", Tag: "payment_method_change", Description: "user change payment method"}
	TopicUserMetricUpdate                 = redismq.MQTopicEnum{Topic: "unibee_user", Tag: "metric_update", Description: "user metric update"}
	TopicUserAccountCreate                = redismq.MQTopicEnum{Topic: "unibee_user", Tag: "account_create", Description: "user account create"}
	TopicUserAccountUpdate                = redismq.MQTopicEnum{Topic: "unibee_user", Tag: "account_update", Description: "user account update"}
	TopicPaymentCreated                   = redismq.MQTopicEnum{Topic: "unibee_payment", Tag: "payment_created", Description: "payment created"}
	TopicPaymentChecker                   = redismq.MQTopicEnum{Topic: "unibee_payment", Tag: "payment_checker", Description: "payment status checker"}
	TopicPaymentCancel                    = redismq.MQTopicEnum{Topic: "unibee_payment", Tag: "payment_cancelled", Description: "payment cancelled"}
	TopicPaymentAuthorized                = redismq.MQTopicEnum{Topic: "unibee_payment", Tag: "payment_authorized", Description: "payment authorized"}
	TopicPaymentSuccess                   = redismq.MQTopicEnum{Topic: "unibee_payment", Tag: "payment_success", Description: "payment success"}
	TopicRefundCreated                    = redismq.MQTopicEnum{Topic: "unibee_refund", Tag: "refund_created", Description: "refund created"}
	TopicRefundChecker                    = redismq.MQTopicEnum{Topic: "unibee_refund", Tag: "refund_checker", Description: "refund status checker"}
	TopicRefundSuccess                    = redismq.MQTopicEnum{Topic: "unibee_refund", Tag: "refund_success", Description: "refund success"}
	TopicRefundFailed                     = redismq.MQTopicEnum{Topic: "unibee_refund", Tag: "refund_failed", Description: "refund success"}
	TopicSubscriptionCancel               = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_cancelled", Description: "subscription cancelled"}
	TopicSubscriptionExpire               = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_expired", Description: "subscription expired"}
	TopicSubscriptionFailed               = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_failed", Description: "subscription failed"}
	TopicSubscriptionCreate               = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_created", Description: "subscription created"}
	TopicSubscriptionUpdate               = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_update", Description: "subscription updated"}
	TopicSubscriptionActive               = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_active", Description: "subscription active"}
	TopicSubscriptionCreatePaymentCheck   = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_create_payment_check", Description: "subscription create payment check"}
	TopicSubscriptionIncomplete           = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_incomplete", Description: "subscription incomplete"}
	TopicSubscriptionPaymentSuccess       = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_payment_success", Description: "subscription payment success"}
	TopicSubscriptionAutoRenewSuccess     = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_auto_renew_success", Description: "subscription auto renew success"}
	TopicSubscriptionAutoRenewFailure     = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_auto_renew_failure", Description: "subscription auto renew failure"}
	TopicMerchantWebhook                  = redismq.MQTopicEnum{Topic: "unibee_merchant_webhook", Tag: "webhook", Description: "merchant webhook"}
	TopicInternalWebhook                  = redismq.MQTopicEnum{Topic: "unibee_internal_webhook", Tag: "webhook", Description: "internal webhook"}
	TopicInvoiceCreated                   = redismq.MQTopicEnum{Topic: "unibee_invoice", Tag: "invoice_created", Description: "invoice created"}
	TopicInvoiceProcessed                 = redismq.MQTopicEnum{Topic: "unibee_invoice", Tag: "invoice_processed", Description: "invoice processed"}
	TopicInvoicePaid                      = redismq.MQTopicEnum{Topic: "unibee_invoice", Tag: "invoice_paid", Description: "invoice paid"}
	TopicInvoiceCancelled                 = redismq.MQTopicEnum{Topic: "unibee_invoice", Tag: "invoice_cancelled", Description: "invoice cancelled"}
	TopicInvoiceFailed                    = redismq.MQTopicEnum{Topic: "unibee_invoice", Tag: "invoice_failed", Description: "invoice failed"}
	TopicInvoiceReversed                  = redismq.MQTopicEnum{Topic: "unibee_invoice", Tag: "invoice_reversed", Description: "invoice reversed"}
	TopicSubscriptionPendingUpdateCancel  = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_pending_update_cancelled", Description: "subscription pending update cancelled"}
	TopicSubscriptionPendingUpdateSuccess = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_pending_update_success", Description: "subscription pending update success"}
	TopicSubscriptionPendingUpdateCreate  = redismq.MQTopicEnum{Topic: "unibee_subscription", Tag: "subscription_pending_update_created", Description: "subscription pending update created"}
	TopicCreditRechargeSuccess            = redismq.MQTopicEnum{Topic: "unibee_credit", Tag: "recharge_success", Description: "recharge success"}
	TopicCreditPaymentSuccess             = redismq.MQTopicEnum{Topic: "unibee_credit", Tag: "payment_success", Description: "payment success"}
	TopicCreditRefundSuccess              = redismq.MQTopicEnum{Topic: "unibee_credit", Tag: "refund_success", Description: "refund success"}
	TopicMerchantCreatedWebhook           = redismq.MQTopicEnum{Topic: "unibee_merchant", Tag: "merchant_create", Description: "merchant create"}
	TopicMerchantMemberCreatedWebhook     = redismq.MQTopicEnum{Topic: "unibee_member", Tag: "member_create", Description: "member create"}
)

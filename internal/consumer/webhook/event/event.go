package event

type MerchantWebhookEvent string

const (
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CREATED   = "subscription.created"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED   = "subscription.updated"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CANCELLED = "subscription.cancelled"
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_EXPIRED   = "subscription.expired"
	MERCHANT_WEBHOOK_TAG_USER_METRIC_UPDATED    = "user.metric.update"
)

var ListeningEventList = []string{
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CREATED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_CANCELLED,
	MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_EXPIRED,
	MERCHANT_WEBHOOK_TAG_USER_METRIC_UPDATED,
}

func WebhookEventInListeningEvents(target MerchantWebhookEvent) bool {
	if len(target) <= 0 {
		return false
	}
	for _, event := range ListeningEventList {
		if event == string(target) {
			return true
		}
	}
	return false
}

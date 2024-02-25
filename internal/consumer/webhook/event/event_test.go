package event

import (
	"fmt"
	"testing"
)

func TestEvent(t *testing.T) {
	fmt.Println(WebhookEventInListeningEvents(MERCHANT_WEBHOOK_TAG_SUBSCRIPTION_UPDATED))
}

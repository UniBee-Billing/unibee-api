package event

import (
	"fmt"
	"testing"
)

func TestEvent(t *testing.T) {
	fmt.Println(WebhookEventInListeningEvents(UNIBEE_WEBHOOK_EVENT_SUBSCRIPTION_UPDATED))
}

package notify

import "unibee/internal/logic/gateway/api/alipay/api/model"

type AlipaySubscriptionNotify struct {
	AlipayNotify
	SubscriptionRequestId        string                             `json:"subscriptionRequestId,omitempty"`
	SubscriptionId               string                             `json:"subscriptionId,omitempty"`
	SubscriptionStatus           model.SubscriptionStatus           `json:"subscriptionStatus,omitempty"`
	SubscriptionNotificationType model.SubscriptionNotificationType `json:"subscriptionNotificationType,omitempty"`
	SubscriptionStartTime        string                             `json:"subscriptionStartTime,omitempty"`
	SubscriptionEndTime          string                             `json:"subscriptionEndTime,omitempty"`
	PeriodRule                   *model.PeriodRule                  `json:"periodRule,omitempty"`
}

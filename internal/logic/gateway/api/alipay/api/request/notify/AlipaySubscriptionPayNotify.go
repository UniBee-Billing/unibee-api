package notify

type AlipaySubscriptionPayNotify struct {
	AlipayNotify
	SubscriptionRequestId string `json:"subscriptionRequestId,omitempty"`
	SubscriptionId        string `json:"subscriptionId,omitempty"`
	PeriodStartTime       string `json:"periodStartTime,omitempty"`
	PeriodEndTime         string `json:"periodEndTime,omitempty"`
	PhaseNo               string `json:"phaseNo,omitempty"`
}

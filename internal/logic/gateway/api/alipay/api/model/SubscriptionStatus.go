package model

type SubscriptionStatus string

const (
	SubscriptionStatus_ACTIVE     SubscriptionStatus = "ACTIVE"
	SubscriptionStatus_TERMINATED SubscriptionStatus = "TERMINATED"
)

type SubscriptionNotificationType string

const (
	SubscriptionNotificationType_CREATE    SubscriptionNotificationType = "CREATE"
	SubscriptionNotificationType_CHANGE    SubscriptionNotificationType = "CHANGE"
	SubscriptionNotificationType_CANCEL    SubscriptionNotificationType = "CANCEL"
	SubscriptionNotificationType_TERMINATE SubscriptionNotificationType = "TERMINATE"
)

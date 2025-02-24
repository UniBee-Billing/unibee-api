package model

type DisputeNotificationType string

const (
	DisputeNotificationType_DISPUTE_CREATED   DisputeNotificationType = "DISPUTE_CREATED"
	DisputeNotificationType_DISPUTE_JUDGED    DisputeNotificationType = "DISPUTE_JUDGED"
	DisputeNotificationType_DISPUTE_CANCELLED DisputeNotificationType = "DISPUTE_CANCELLED"
	DisputeNotificationType_DEFENSE_SUPPLIED  DisputeNotificationType = "DEFENSE_SUPPLIED"
	DisputeNotificationType_DEFENSE_DUE_ALERT DisputeNotificationType = "DEFENSE_DUE_ALERT"
	DisputeNotificationType_DISPUTE_ACCEPTED  DisputeNotificationType = "DISPUTE_ACCEPTED"
	RDR_RESOLVED                              DisputeNotificationType = "RDR_RESOLVED"
)

type DisputeJudgedResult string

const (
	DisputeJudgedResult_ACCEPT_BY_CUSTOMER DisputeJudgedResult = "ACCEPT_BY_CUSTOMER"
	DisputeJudgedResult_ACCEPT_BY_MERCHANT DisputeJudgedResult = "ACCEPT_BY_MERCHANT"
)

type DisputeAcceptReasonType string

const (
	MERCHANT_ACCEPTED          DisputeAcceptReasonType = "MERCHANT_ACCEPTED"
	TIMEOUT                    DisputeAcceptReasonType = "TIMEOUT"
	MANUAL_PROCESSING_ACCEPTED DisputeAcceptReasonType = "MANUAL_PROCESSING_ACCEPTED"
)

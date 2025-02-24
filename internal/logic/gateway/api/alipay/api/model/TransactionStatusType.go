package model

type TransactionStatusType string

const (
	TransactionStatusType_SUCCESS    TransactionStatusType = "SUCCESS"
	TransactionStatusType_FAIL       TransactionStatusType = "FAIL"
	TransactionStatusType_PROCESSING TransactionStatusType = "PROCESSING"
	TransactionStatusType_PENDING    TransactionStatusType = "PENDING"
	TransactionStatusType_CANCELLED  TransactionStatusType = "CANCELLED"
)

package model

import "unibee/internal/logic/gateway/api/alipay/api/response"

type TransactionType string

const (
	TransactionType_PAYMENT       TransactionType = "PAYMENT"
	TransactionType_REFUND        TransactionType = "REFUND"
	TransactionType_CANCEL        TransactionType = "CANCEL"
	TransactionType_CAPTURE       TransactionType = "CAPTURE"
	TransactionType_VOID          TransactionType = "VOID"
	TransactionType_AUTHORIZATION TransactionType = "AUTHORIZATION"
)

type Transaction struct {
	TransactionResult    response.Result       `json:"transactionResult,omitempty"`
	TransactionId        string                `json:"transactionId,omitempty"`
	TransactionType      TransactionType       `json:"transactionType,omitempty"`
	TransactionStatus    TransactionStatusType `json:"transactionStatus,omitempty"`
	TransactionAmount    Amount                `json:"transactionAmount,omitempty"`
	TransactionRequestId string                `json:"transactionRequestId,omitempty"`
	TransactionTime      string                `json:"transactionTime,omitempty"`

	AcquirerInfo AcquirerInfo `json:"acquirerInfo,omitempty"`
}

package model

type TransferToDetail struct {
	TransferToMethod       *PaymentMethod `json:"transferToMethod,omitempty"`
	TransferToCurrency     string         `json:"transferToCurrency,omitempty"`
	FeeAmount              *Amount        `json:"feeAmount,omitempty"`
	ActualTransferToAmount *Amount        `json:"actualTransferToAmount,omitempty"`
	PurposeCode            string         `json:"purposeCode,omitempty"`
	TransferNotifyUrl      string         `json:"transferNotifyUrl,omitempty"`
	TransferRemark         string         `json:"transferRemark,omitempty"`
}

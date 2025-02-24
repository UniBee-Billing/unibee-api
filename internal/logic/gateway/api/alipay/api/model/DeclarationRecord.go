package model

type DeclarationRecord struct {
	DeclarationRequestId         string              `json:"declarationRequestId,omitempty"`
	CustomsPaymentId             string              `json:"customsPaymentId,omitempty"`
	CustomsOrderId               string              `json:"customsOrderId,omitempty"`
	Customs                      CustomsInfo         `json:"customs,omitempty"`
	MerchantCustomsInfo          MerchantCustomsInfo `json:"merchantCustomsInfo,omitempty"`
	DeclarationAmount            Amount              `json:"declarationAmount,omitempty"`
	SplitOrder                   bool                `json:"splitOrder,omitempty"`
	DeclarationRequestStatus     string              `json:"declarationRequestStatus,omitempty"`
	LastModifiedTime             string              `json:"lastModifiedTime,omitempty"`
	CustomsDeclarationResultCode string              `json:"customsDeclarationResultCode,omitempty"`
	CustomsDeclarationResultDesc string              `json:"customsDeclarationResultDesc,omitempty"`
	CustomsDeclarationReturnTime string              `json:"customsDeclarationReturnTime,omitempty"`
}

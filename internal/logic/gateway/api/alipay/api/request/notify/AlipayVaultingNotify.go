package notify

import "unibee/internal/logic/gateway/api/alipay/api/model"

type AlipayVaultingNotify struct {
	AlipayNotify
	VaultingRequestId   string                     `json:"vaultingRequestId,omitempty"`
	PaymentMethodDetail *model.PaymentMethodDetail `json:"paymentMethodDetail,omitempty"`
	VaultingCreateTime  string                     `json:"vaultingCreateTime,omitempty"`
	AcquirerInfo        string                     `json:"acquirerInfo,omitempty"`
}

package responseVaulting

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayVaultingQueryResponse struct {
	response.AlipayResponse
	VaultingRequestId   string                     `json:"vaultingRequestId,omitempty"`
	NormalUrl           string                     `json:"normalUrl,omitempty"`
	SchemeUrl           string                     `json:"schemeUrl,omitempty"`
	ApplinkUrl          string                     `json:"applinkUrl,omitempty"`
	VaultingStatus      string                     `json:"vaultingStatus,omitempty"`
	PaymentMethodDetail *model.PaymentMethodDetail `json:"paymentMethodDetail,omitempty"`
}

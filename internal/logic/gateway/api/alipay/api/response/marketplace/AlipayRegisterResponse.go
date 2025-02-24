package responseMarketplace

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipayRegisterResponse struct {
	response.AlipayResponse
	RegistrationStatus string `json:"registrationStatus,omitempty"`
}

package responseAuth

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipayAuthCreateSessionResponse struct {
	response.AlipayResponse
	PaymentSessionId         string `json:"paymentSessionId,omitempty"`
	PaymentSessionData       string `json:"paymentSessionData,omitempty"`
	PaymentSessionExpiryTime string `json:"paymentSessionExpiryTime,omitempty"`
}

package responsePay

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipayPaymentSessionResponse struct {
	response.AlipayResponse
	PaymentSessionData       string `json:"paymentSessionData"`
	PaymentSessionExpiryTime string `json:"paymentSessionExpiryTime"`
	PaymentSessionId         string `json:"paymentSessionId"`
	NormalUrl                string `json:"normalUrl"`
}

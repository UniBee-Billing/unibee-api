package responsePay

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipayPayCancelResponse struct {
	response.AlipayResponse
	PaymentId        string `json:"paymentId"`
	PaymentRequestId string `json:"paymentRequestId"`
	CancelTime       string `json:"cancelTime"`
}

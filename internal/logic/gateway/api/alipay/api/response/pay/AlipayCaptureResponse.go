package responsePay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayCaptureResponse struct {
	response.AlipayResponse
	CaptureRequestId    string       `json:"captureRequestId"`
	CaptureId           string       `json:"captureId"`
	PaymentId           string       `json:"paymentId"`
	CaptureAmount       model.Amount `json:"captureAmount"`
	CaptureTime         string       `json:"captureTime"`
	AcquirerReferenceNo string       `json:"acquirerReferenceNo"`
}

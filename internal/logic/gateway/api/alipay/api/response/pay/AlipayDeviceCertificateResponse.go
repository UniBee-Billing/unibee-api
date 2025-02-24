package responsePay

import "unibee/internal/logic/gateway/api/alipay/api/response"

type AlipayDeviceCertificateResponse struct {
	response.AlipayResponse
	DeviceCertificate string `json:"deviceCertificate,omitempty"`
}

package pay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responsePay "unibee/internal/logic/gateway/api/alipay/api/response/pay"
)

type AlipayDeviceCertificateRequest struct {
	DevicePublicKey string `json:"devicePublicKey,omitempty"`
	DeviceRequestId string `json:"deviceRequestId,omitempty"`
}

func (alipayCaptureRequest *AlipayDeviceCertificateRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayCaptureRequest, model.CREATE_DEVICE_CERTIFICATE_PATH, &responsePay.AlipayDeviceCertificateResponse{})
}

func NewAlipayDeviceCertificateRequest() (*request.AlipayRequest, *AlipayDeviceCertificateRequest) {
	alipayDeviceCertificateRequest := &AlipayDeviceCertificateRequest{}
	alipayRequest := request.NewAlipayRequest(alipayDeviceCertificateRequest, model.CREATE_DEVICE_CERTIFICATE_PATH, &responsePay.AlipayDeviceCertificateResponse{})
	return alipayRequest, alipayDeviceCertificateRequest
}

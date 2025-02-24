package pay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responsePay "unibee/internal/logic/gateway/api/alipay/api/response/pay"
)

type AlipayCaptureRequest struct {
	CaptureRequestId string        `json:"captureRequestId,omitempty"`
	PaymentId        string        `json:"paymentId,omitempty"`
	CaptureAmount    *model.Amount `json:"captureAmount,omitempty"`
	IsLastCapture    bool          `json:"isLastCapture,omitempty"`
}

func (alipayCaptureRequest *AlipayCaptureRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayCaptureRequest, model.CAPTURE_PATH, &responsePay.AlipayCaptureResponse{})
}

func NewAlipayCaptureRequest() (*request.AlipayRequest, *AlipayCaptureRequest) {
	alipayCaptureRequest := &AlipayCaptureRequest{}
	alipayRequest := request.NewAlipayRequest(alipayCaptureRequest, model.CAPTURE_PATH, &responsePay.AlipayCaptureResponse{})
	return alipayRequest, alipayCaptureRequest
}

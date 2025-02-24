package pay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responsePay "unibee/internal/logic/gateway/api/alipay/api/response/pay"
)

type AlipaySyncArrearRequest struct {
	PaymentId        string `json:"paymentId,omitempty"`
	PaymentRequestId string `json:"paymentRequestId,omitempty"`
}

func (alipayCaptureRequest *AlipaySyncArrearRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayCaptureRequest, model.SYNC_ARREAR_PATH, &responsePay.AlipaySyncArrearResponse{})
}

func NewAlipaySyncArrearRequest() (*request.AlipayRequest, *AlipaySyncArrearRequest) {
	alipaySyncArrearRequest := &AlipaySyncArrearRequest{}
	alipayRequest := request.NewAlipayRequest(alipaySyncArrearRequest, model.SYNC_ARREAR_PATH, &responsePay.AlipaySyncArrearResponse{})
	return alipayRequest, alipaySyncArrearRequest
}

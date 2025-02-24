package pay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	"unibee/internal/logic/gateway/api/alipay/api/response/pay"
)

type AlipayPayQueryRequest struct {
	PaymentRequestId  string `json:"paymentRequestId,omitempty"`
	PaymentId         string `json:"paymentId,omitempty"`
	MerchantAccountId string `json:"MerchantAccountId,omitempty"`
}

func (alipayPayQueryRequest *AlipayPayQueryRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayPayQueryRequest, model.INQUIRY_PAYMENT_PATH, &responsePay.AlipayPayQueryResponse{})
}

func NewAlipayPayQueryRequest() (*request.AlipayRequest, *AlipayPayQueryRequest) {
	alipayPayQueryRequest := &AlipayPayQueryRequest{}
	alipayRequest := request.NewAlipayRequest(alipayPayQueryRequest, model.INQUIRY_PAYMENT_PATH, &responsePay.AlipayPayQueryResponse{})
	return alipayRequest, alipayPayQueryRequest
}

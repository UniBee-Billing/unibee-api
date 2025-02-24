package pay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responsePay "unibee/internal/logic/gateway/api/alipay/api/response/pay"
)

type AlipayPayCancelRequest struct {
	PaymentId         string `json:"paymentId,omitempty"`
	PaymentRequestId  string `json:"paymentRequestId,omitempty"`
	MerchantAccountId string `json:"merchantAccountId,omitempty"`
}

func NewAlipayPayCancelRequest() (*request.AlipayRequest, *AlipayPayCancelRequest) {
	alipayPayCancelRequest := &AlipayPayCancelRequest{}
	alipayRequest := request.NewAlipayRequest(alipayPayCancelRequest, model.CANCEL_PATH, &responsePay.AlipayPayCancelResponse{})
	return alipayRequest, alipayPayCancelRequest
}

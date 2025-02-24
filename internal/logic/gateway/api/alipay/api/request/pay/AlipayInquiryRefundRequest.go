package pay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responsePay "unibee/internal/logic/gateway/api/alipay/api/response/pay"
)

type AlipayInquiryRefundRequest struct {
	RefundRequestId   string `json:"refundRequestId,omitempty"`
	RefundId          string `json:"refundId,omitempty"`
	MerchantAccountId string `json:"merchantAccountId,omitempty"`
}

func (alipayInquiryRefundRequest *AlipayInquiryRefundRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayInquiryRefundRequest, model.INQUIRY_REFUND_PATH, &responsePay.AlipayInquiryRefundResponse{})
}

func NewAlipayInquiryRefundRequest() (*request.AlipayRequest, *AlipayInquiryRefundRequest) {
	alipayInquiryRefundRequest := &AlipayInquiryRefundRequest{}
	alipayRequest := request.NewAlipayRequest(alipayInquiryRefundRequest, model.INQUIRY_REFUND_PATH, &responsePay.AlipayInquiryRefundResponse{})
	return alipayRequest, alipayInquiryRefundRequest
}

package responsePay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayInquiryRefundResponse struct {
	response.AlipayResponse
	RefundId              string                      `json:"refundId"`
	RefundRequestId       string                      `json:"refundRequestId"`
	RefundAmount          model.Amount                `json:"refundAmount"`
	RefundStatus          model.TransactionStatusType `json:"refundStatus"`
	RefundTime            string                      `json:"refundTime"`
	GrossSettlementAmount model.Amount                `json:"grossSettlementAmount"`
	SettlementQuote       model.Quote                 `json:"settlementQuote"`
	AcquirerInfo          model.AcquirerInfo          `json:"acquirerInfo"`
}

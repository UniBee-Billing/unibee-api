package responsePay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/response"
)

type AlipayPayConsultResponse struct {
	response.AlipayResponse
	PaymentOptions     []model.PaymentOption     `json:"paymentOptions"`
	PaymentMethodInfos []model.PaymentMethodInfo `json:"paymentMethodInfos"`
	ExtendInfo         string                    `json:"extendInfo"`
}

package pay

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responsePay "unibee/internal/logic/gateway/api/alipay/api/response/pay"
)

type AlipayPayConsultRequest struct {
	ProductCode                 model.ProductCodeType     `json:"productCode,omitempty"`
	PaymentAmount               *model.Amount             `json:"paymentAmount,omitempty"`
	MerchantRegion              string                    `json:"merchantRegion,omitempty"`
	AllowedPaymentMethodRegions []string                  `json:"allowedPaymentMethodRegions,omitempty"`
	AllowedPaymentMethods       []string                  `json:"allowedPaymentMethods,omitempty"`
	BlockedPaymentMethods       []string                  `json:"blockedPaymentMethods,omitempty"`
	Region                      string                    `json:"region,omitempty"`
	CustomerId                  string                    `json:"customerId,omitempty"`
	ReferenceUserId             string                    `json:"referenceUserId,omitempty"`
	Env                         *model.Env                `json:"env,omitempty"`
	ExtendInfo                  string                    `json:"extendInfo,omitempty"`
	UserRegion                  string                    `json:"userRegion,omitempty"`
	PaymentFactor               *model.PaymentFactor      `json:"paymentFactor,omitempty"`
	SettlementStrategy          *model.SettlementStrategy `json:"settlementStrategy,omitempty"`
	Merchant                    *model.Merchant           `json:"merchant,omitempty"`
	AllowedPspRegions           []string                  `json:"allowedPspRegions,omitempty"`
	Buyer                       *model.Buyer              `json:"buyer,omitempty"`
	MerchantAccountId           string                    `json:"merchantAccountId,omitempty"`
}

func (alipayPayConsultRequest *AlipayPayConsultRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayPayConsultRequest, model.CONSULT_PAYMENT_PATH, &responsePay.AlipayPayConsultResponse{})
}

func NewAlipayPayConsultRequest() (*request.AlipayRequest, *AlipayPayConsultRequest) {
	alipayPayConsultRequest := &AlipayPayConsultRequest{}
	alipayRequest := request.NewAlipayRequest(alipayPayConsultRequest, model.CONSULT_PAYMENT_PATH, &responsePay.AlipayPayConsultResponse{})
	return alipayRequest, alipayPayConsultRequest
}

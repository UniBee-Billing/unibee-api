package auth

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseAuth "unibee/internal/logic/gateway/api/alipay/api/response/auth"
)

type AlipayAuthCreateSessionRequest struct {
	ProductCode        model.ProductCodeType `json:"productCode,omitempty"`
	AgreementInfo      *model.AgreementInfo  `json:"agreementInfo,omitempty"`
	Scopes             []model.ScopeType     `json:"scopes,omitempty"`
	PaymentMethod      *model.PaymentMethod  `json:"paymentMethod,omitempty"`
	PaymentRedirectUrl string                `json:"paymentRedirectUrl,omitempty"`
}

func (alipayAuthCreateSessionRequest *AlipayAuthCreateSessionRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayAuthCreateSessionRequest, model.CREATE_SESSION_PATH, &responseAuth.AlipayAuthCreateSessionResponse{})
}

func NewAlipayAuthCreateSessionRequest() (*request.AlipayRequest, *AlipayAuthCreateSessionRequest) {
	alipayAuthCreateSessionRequest := &AlipayAuthCreateSessionRequest{}
	alipayRequest := request.NewAlipayRequest(alipayAuthCreateSessionRequest, model.CREATE_SESSION_PATH, &responseAuth.AlipayAuthCreateSessionResponse{})
	return alipayRequest, alipayAuthCreateSessionRequest
}

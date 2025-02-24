package auth

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseAuth "unibee/internal/logic/gateway/api/alipay/api/response/auth"
)

type AlipayAuthConsultRequest struct {
	CustomerBelongsTo model.CustomerBelongsTo `json:"customerBelongsTo,omitempty"`
	AuthClientId      string                  `json:"authClientId,omitempty"`
	AuthRedirectUrl   string                  `json:"authRedirectUrl,omitempty"`
	Scopes            []model.ScopeType       `json:"scopes,omitempty"`
	AuthState         string                  `json:"authState,omitempty"`
	TerminalType      model.TerminalType      `json:"terminalType,omitempty"`
	OsType            model.OsType            `json:"osType,omitempty"`
	OsVersion         string                  `json:"osVersion,omitempty"`
	ExtendInfo        string                  `json:"extendInfo,omitempty"`
	MerchantRegion    string                  `json:"merchantRegion,omitempty"`
	RecurringPayment  bool                    `json:"recurringPayment,omitempty"`
	GrantType         string                  `json:"grantType,omitempty"`
	AuthMetaData      *model.AuthMetaData     `json:"authMetaData,omitempty"`
}

func (alipayAuthConsultRequest *AlipayAuthConsultRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayAuthConsultRequest, model.AUTH_CONSULT_PATH, &responseAuth.AlipayAuthConsultResponse{})
}

func NewAlipayAuthConsultRequest() (*request.AlipayRequest, *AlipayAuthConsultRequest) {
	alipayAuthConsultRequest := &AlipayAuthConsultRequest{}
	alipayRequest := request.NewAlipayRequest(alipayAuthConsultRequest, model.AUTH_CONSULT_PATH, &responseAuth.AlipayAuthConsultResponse{})
	return alipayRequest, alipayAuthConsultRequest
}

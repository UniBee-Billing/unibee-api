package auth

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseAuth "unibee/internal/logic/gateway/api/alipay/api/response/auth"
)

type AlipayAuthApplyTokenRequest struct {
	GrantType         model.GrantType         `json:"grantType,omitempty"`
	CustomerBelongsTo model.CustomerBelongsTo `json:"customerBelongsTo,omitempty"`
	AuthCode          string                  `json:"authCode,omitempty"`
	RefreshToken      string                  `json:"refreshToken,omitempty"`
	ExtendInfo        string                  `json:"extendInfo,omitempty"`
	MerchantRegion    string                  `json:"merchantRegion,omitempty"`
}

func (alipayAuthApplyTokenRequest *AlipayAuthApplyTokenRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayAuthApplyTokenRequest, model.AUTH_APPLY_TOKEN_PATH, &responseAuth.AlipayAuthApplyTokenResponse{})
}

func NewAlipayAuthApplyTokenRequest() (*request.AlipayRequest, *AlipayAuthApplyTokenRequest) {
	alipayAuthApplyTokenRequest := &AlipayAuthApplyTokenRequest{}
	alipayRequest := request.NewAlipayRequest(alipayAuthApplyTokenRequest, model.AUTH_APPLY_TOKEN_PATH, &responseAuth.AlipayAuthApplyTokenResponse{})
	return alipayRequest, alipayAuthApplyTokenRequest
}

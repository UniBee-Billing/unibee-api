package auth

import (
	"unibee/internal/logic/gateway/api/alipay/api/model"
	"unibee/internal/logic/gateway/api/alipay/api/request"
	responseAuth "unibee/internal/logic/gateway/api/alipay/api/response/auth"
)

type AlipayAuthRevokeTokenRequest struct {
	AccessToken string `json:"accessToken,omitempty"`
	ExtendInfo  string `json:"extendInfo,omitempty"`
}

func (alipayAuthRevokeTokenRequest *AlipayAuthRevokeTokenRequest) NewRequest() *request.AlipayRequest {
	return request.NewAlipayRequest(&alipayAuthRevokeTokenRequest, model.AUTH_REVOKE_PATH, &responseAuth.AlipayAuthRevokeTokenResponse{})
}

func NewAlipayAuthRevokeTokenRequest() (*request.AlipayRequest, *AlipayAuthRevokeTokenRequest) {
	alipayAuthRevokeTokenRequest := &AlipayAuthRevokeTokenRequest{}
	alipayRequest := request.NewAlipayRequest(alipayAuthRevokeTokenRequest, model.AUTH_REVOKE_PATH, &responseAuth.AlipayAuthRevokeTokenResponse{})
	return alipayRequest, alipayAuthRevokeTokenRequest
}
